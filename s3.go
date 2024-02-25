package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/majesticbeast/terraspoof/internal/database"
	"net/http"
	"time"
)

// createS3Bucket will add an entry to the s3 table.
func (a *ApiServer) createS3Bucket(w http.ResponseWriter, r *http.Request) error {
	var params struct {
		Bucket string `json:"bucket"`
		Tags   string `json:"tags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return err
	}

	result, err := a.db.CreateS3(context.Background(), database.CreateS3Params{
		ID:               uuid.New(),
		Bucket:           params.Bucket,
		Tags:             params.Tags,
		BucketDomainName: params.Bucket + ".s3.majestic-cloud.com",
		CreatedAt:        time.Now().UTC(),
	})
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"s3_bucket_key\"" {
			err = ErrBucketExists
		}
		return err
	}

	respondWithJSON(w, http.StatusCreated, result)
	a.slog.Info("S3 bucket created", "bucket", "id", result.Bucket, result.ID.String())
	return nil
}

// deleteS3Bucket will delete an entry from the s3 table.
func (a *ApiServer) deleteS3Bucket(w http.ResponseWriter, r *http.Request) error {
	var params struct {
		Bucket string `json:"bucket"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return err
	}

	err := a.db.DeleteS3(context.Background(), params.Bucket)
	if err != nil {
		return err
	}

	respondWithJSON(w, http.StatusNoContent, nil)
	a.slog.Info("S3 bucket deleted", "bucket", params)
	return nil
}

// getS3Bucket will get the information of a specific s3 bucket.
func (a *ApiServer) getS3Bucket(w http.ResponseWriter, r *http.Request) error {
	var params struct {
		Bucket string `json:"bucket"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return err
	}

	result, err := a.db.GetS3(context.Background(), params.Bucket)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "S3 bucket not found")
		return err
	}

	respondWithJSON(w, http.StatusOK, result)
	a.slog.Info("Listed all S3 buckets")
	return nil
}
