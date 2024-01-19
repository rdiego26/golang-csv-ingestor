package main

import (
	"database/sql"
	"encoding/json"
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
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("Please set the environment variable DATABASE_URL")
	}

	log.Println("Connecting with database...")
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/users", getUsers(db)).Methods("GET")

	//http.ListenAndServeTLS()
	log.Println("API server running on port: ", s.listenAddr)
	log.Fatal(http.ListenAndServe(s.listenAddr, jsonContentTypeMiddleware(router)))
}

func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		// Convert the url.Values map to a map[string]string
		urlParams := make(map[string]string)
		for key, values := range queryParams {
			if len(values) > 0 {
				urlParams[key] = values[0]
			}
		}

		_ = filterParams(urlParams)

		// TODO: build query based on params
		rows, err := db.Query("SELECT id, first_name, last_name, email, parent_id FROM users")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.ParentId); err != nil {
				log.Fatal(err)
			}
			users = append(users, u)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(users)
	}
}

func filterParams(receivedParams map[string]string) map[string]string {
	result := make(map[string]string)
	allowedParams := []string{"first_name", "last_name", "email", "parent_id"}

	for key, value := range receivedParams {
		for _, allowedKey := range allowedParams {
			if key == allowedKey {
				result[key] = value
				break
			}
		}
	}

	log.Println("Received parameters:", receivedParams)
	log.Println("Filtered parameters:", result)

	return result
}
