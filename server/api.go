package server

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type APIServer struct {
	listenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) Run() {
	// Create a new logger instance.
	logger := &Logger{}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		logger.Log(Fatal, "Please set the environment variable DATABASE_URL")
	}

	logger.Log(Info, "Connecting with database...")
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/users", getUsers(db)).Methods("GET")

	//http.ListenAndServeTLS()
	logger.Log(Info, "API server running on port: "+s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, jsonContentTypeMiddleware(router)))
}
