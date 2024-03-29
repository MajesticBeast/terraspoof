package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/majesticbeast/terraspoof/internal/database"
	"log/slog"
	"net/http"
)

var (
	ErrBucketExists = errors.New("bucket already exists")
)

type ApiServer struct {
	port string
	slog *slog.Logger
	db   *database.Queries
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

// makeHTTPHandler will create an HTTP handler for the specified API function. This is used to handle errors returned
// from the HTTP handlers.
func (a *ApiServer) makeHTTPHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			a.slog.Error(err.Error())
		}
	}
}

// NewApiServer will create a new API server with the specified port, logger, and database connection.
func NewApiServer(port string, slog *slog.Logger, db *database.Queries) *ApiServer {
	return &ApiServer{
		port: port,
		slog: slog,
		db:   db,
	}
}

// Start will start the API server and listen on the specified port.
func (a *ApiServer) Start() error {

	// Create routers and subrouters
	router := mux.NewRouter()
	s3Router := router.PathPrefix("/api/v1/s3").Subrouter()
	userRouter := router.PathPrefix("/api/v1/users").Subrouter()

	// Register main routes
	router.HandleFunc("/api/v1/health", a.HealthCheck).Methods("GET")

	// Register s3 subroutes
	s3Router.HandleFunc("/create", a.makeHTTPHandler(a.createS3Bucket)).Methods("POST")
	s3Router.HandleFunc("/delete", a.makeHTTPHandler(a.deleteS3Bucket)).Methods("DELETE")
	s3Router.HandleFunc("/get", a.makeHTTPHandler(a.getS3Bucket)).Methods("GET")

	// Register user subroutes
	userRouter.HandleFunc("/create", a.makeHTTPHandler(a.createUser)).Methods("POST")
	userRouter.HandleFunc("/delete", a.makeHTTPHandler(a.deleteUser)).Methods("DELETE")
	userRouter.HandleFunc("/get", a.makeHTTPHandler(a.getUser)).Methods("GET")
	userRouter.HandleFunc("/login", a.makeHTTPHandler(a.login)).Methods("POST")

	// Start the server
	a.slog.Info("Server listening on " + a.port)
	return http.ListenAndServe(a.port, router)
}

// HealthCheck will return a 200 OK response.
func (a *ApiServer) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

// respondWithJSON will respond with a JSON payload and a non-error status code
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondWithError will respond with a JSON payload and an error status code
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
