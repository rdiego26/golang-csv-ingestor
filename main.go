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

	server := NewAPIServer(":" + appPort)
	server.Run()
}
