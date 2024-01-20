package main

import (
	"database/sql"
	"log"
	"os"
)

func main() {
	appPort := os.Getenv("APP_PORT")

	if appPort == "" {
		log.Fatal("Please set the environment variable APP_PORT")
	}

	go cleanTableBeforeImport()
	go importData()

	server := NewAPIServer(":" + appPort)
	server.Run()
}

func cleanTableBeforeImport() {
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

	logger.Log(Info, "Truncating users table...")
	truncateUsers(db)
}

func importData() {
	filePath := "./csv/users.csv"
	fetchedContent := Read(filePath)
	parsedUsers := ProcessContacts(fetchedContent)

	for _, user := range parsedUsers {
		sendMessage(user)
	}
}
