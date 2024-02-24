package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/majesticbeast/terraspoof/internal/database"
	"net/http"
	"time"
)

// This endpoint will add an entry to the s3 table.
func (a *ApiServer) createS3Bucket(w http.ResponseWriter, r *http.Request) error {
	var params database.CreateS3Params
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return err
	}

	params.CreatedAt = time.Now()
	params.ID = uuid.New()
	params.BucketDomainName = params.Bucket + ".s3.majestic-cloud.com"

	result, err := a.db.CreateS3(context.Background(), params)
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

// This endpoint will delete an entry from the s3 table.
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
