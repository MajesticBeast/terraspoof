package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/majesticbeast/terraspoof/internal/database"
	"log/slog"
	"os"
)

const (
	port     = ":8080"
	dbDriver = "postgres"
)

func main() {
	// Grab environment variables
	godotenv.Load(".env")
	PGCONN := os.Getenv("PGCONN")

	// Setup structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Setup database connection
	db, err := sql.Open(dbDriver, PGCONN)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	dbQueries := database.New(db)

	// Start the API server
	router := NewApiServer(port, logger, dbQueries)
	if err := router.Start(); err != nil {
		logger.Error(err.Error())
	}
}
