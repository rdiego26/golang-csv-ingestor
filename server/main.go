package server

import (
	"fmt"
	"log"
	"os"
)

func main() {
	appPort := os.Getenv("APP_PORT")

	if appPort == "" {
		log.Fatal("Please set the environment variable APP_PORT")
	}

	cleanTableBeforeImport()
	go processAndPublishData()
	go startConsumer()

	server := NewAPIServer(":" + appPort)
	server.Run()
}

func cleanTableBeforeImport() {
	logger := &Logger{}
	db := getConnection()

	logger.Log(Info, "Truncating users table...")
	truncateUsers(db)
}

func processAndPublishData() {
	filePath := "./csv/users.csv"
	fetchedContent := Read(filePath)
	parsedUsers := ProcessContacts(fetchedContent)

	logger := &Logger{}
	logger.Log(Debug, fmt.Sprintf("Bulking %d users", len(parsedUsers)))
	bulkInsertUsers(getConnection(), parsedUsers)

	for _, user := range parsedUsers {
		sendMessage(user)
	}

}

func startConsumer() {
	logger := &Logger{}
	logger.Log(Debug, "Consumer was started!")
	launchConsumer()
}
