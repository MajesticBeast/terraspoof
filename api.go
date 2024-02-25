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

func (a *ApiServer) makeHTTPHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			//var code int
			//switch {
			//case errors.Is(err, ErrBucketExists):
			//	code = http.StatusConflict
			//}
			//http.Error(w, err.Error(), code)
			a.slog.Error(err.Error())
		}
	}
}

func NewApiServer(port string, slog *slog.Logger, db *database.Queries) *ApiServer {
	return &ApiServer{
		port: port,
		slog: slog,
		db:   db,
	}
}

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

	// Start the server
	a.slog.Info("Server listening on " + a.port)
	return http.ListenAndServe(a.port, router)
}

func (a *ApiServer) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
