package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// login will authenticate the user and return the user object.
func (a *ApiServer) login(w http.ResponseWriter, r *http.Request) error {
	var params struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return err
	}

	user, err := a.db.GetUserByName(context.Background(), params.Name)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return err
	}

	if !doPasswordsMatch(user.Password, params.Password) {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return err
	}

	userResponse := struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}{
		ID:   user.ID,
		Name: user.Name,
	}

	respondWithJSON(w, http.StatusOK, userResponse)
	a.slog.Info("User logged in", "name", user.Name)
	return nil
}

// generateAPIKey will generate a random 32 byte string and encode it to base64.
func generateAPIKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// hashPassword will hash the password using bcrypt.
func hashPassword(password string) (string, error) {
	var passwordBytes = []byte(password)
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	return string(hashedPasswordBytes), err
}

// doPasswordsMatch will compare the hashed password with the current password.
func doPasswordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil
}
