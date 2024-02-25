-- name: CreateS3 :one
INSERT INTO s3 (id, bucket, tags, bucket_domain_name, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteS3 :exec
DELETE FROM s3
WHERE bucket = $1;

-- name: GetS3 :one
SELECT * FROM s3
WHERE bucket = $1;