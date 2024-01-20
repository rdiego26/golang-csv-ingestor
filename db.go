package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func truncateUsers(db *sql.DB) {
	query := "TRUNCATE TABLE users"

	result, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
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
