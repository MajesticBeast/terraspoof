package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/majesticbeast/terraspoof/internal/database"
	"net/http"
	"time"
)

// createUser will create a new user in the database.
func (a *ApiServer) createUser(w http.ResponseWriter, r *http.Request) error {
	var params struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return err
	}

	apiKey, err := generateAPIKey()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user: code 1")
		return err
	}

	hashedPassword, err := hashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user: code 2")
		return err
	}

	result, err := a.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		Password:  hashedPassword,
		ApiKey:    apiKey,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user: code 3")
		return err
	}

	respondWithJSON(w, http.StatusCreated, result)
	a.slog.Info("User created", "name", result.Name)
	return nil
}

// deleteUser will delete a user from the database.
func (a *ApiServer) deleteUser(w http.ResponseWriter, r *http.Request) error {
	var params struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return err
	}

	err := a.db.DeleteUser(context.Background(), params.Name)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return err
	}

	respondWithJSON(w, http.StatusNoContent, nil)
	a.slog.Info("User deleted", "name", params.Name)
	return nil
}

// getUser will retrieve a user from the database by name.
func (a *ApiServer) getUser(w http.ResponseWriter, r *http.Request) error {
	var params struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return err
	}

	result, err := a.db.GetUserByName(context.Background(), params.Name)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return err
	}

	var userResponse = struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}{
		ID:   result.ID,
		Name: result.Name,
	}

	respondWithJSON(w, http.StatusOK, userResponse)
	a.slog.Info("User retrieved", "name", params.Name)
	return nil
}
