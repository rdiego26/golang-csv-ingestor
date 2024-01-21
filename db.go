package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func getConnection() *sql.DB {
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

	return db
}

func truncateUsers(db *sql.DB) {
	logger := &Logger{}
	query := "TRUNCATE TABLE users"

	result, err := db.Query(query)
	if err != nil {
		logger.Log(Fatal, fmt.Sprintf("Error while truncate users table: %v", err))
	}
	defer result.Close()
}

func insertUser(db *sql.DB, user User) {
	logger := &Logger{}
	query := "INSERT INTO users(id, first_name, last_name, email, parent_id, created_at, deleted_at, merged_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

	result, err := db.Query(query, user.ID, user.FirstName, user.LastName, user.Email, user.CreatedAt, user.DeletedAt, user.MergedAt)
	if err != nil {
		logger.Log(Error, fmt.Sprintf("Error while add user into users table: %v", err))
	}
	defer result.Close()

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

		mapParams := filterParams(urlParams)
		query := "SELECT id, first_name, last_name, email, parent_id FROM users"
		args := make([]interface{}, 0)
		if len(mapParams) != 0 {
			index := 1
			query = query + " WHERE 1=1 "
			for key, value := range mapParams {
				query += fmt.Sprintf(" AND %s = $%d", key, index)
				args = append(args, value)
				index++
			}
		} else {
			// FIXME: implement pagination
			query = query + " LIMIT 100"
		}

		rows, err := db.Query(query, args...)
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
	logger := &Logger{}

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

	if len(receivedParams) != len(result) {
		logger.Log(Warning, fmt.Sprintf("Received parameters: %s allowed parameters after validation: %s", receivedParams, result))
	}

	return result
}
