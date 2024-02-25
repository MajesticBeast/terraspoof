package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/majesticbeast/terraspoof/internal/database"
	"net/http"
	"time"
)

func (a *ApiServer) createUser(w http.ResponseWriter, r *http.Request) error {
	var params struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return err
	}

	apiKey, err := generateAPIKey()
	if err != nil {
		return err
	}
	result, err := a.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		Password:  params.Password,
		ApiKey:    apiKey,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	respondWithJSON(w, http.StatusCreated, result)
	a.slog.Info("User created", "name", result.Name)
	return nil
}

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

	respondWithJSON(w, http.StatusOK, result)
	a.slog.Info("User retrieved", "name", params.Name)
	return nil
}

func generateAPIKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}
