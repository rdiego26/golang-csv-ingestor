package main

import (
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
	//go startConsumer()

	server := NewAPIServer(":" + appPort)
	server.Run()
}

func cleanTableBeforeImport() {
	logger := &Logger{}
	db := getConnection()

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

func startConsumer() {
	launchConsumer()
}
