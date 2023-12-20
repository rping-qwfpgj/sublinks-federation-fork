package main

import (
	"fmt"
	"sublinks/sublinks-federation/internal/db"
	"sublinks/sublinks-federation/internal/http"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/queue"

	"github.com/joho/godotenv"
)

func main() {
	// bootstrap logger
	logger := log.NewLogger("main")

	// Load connection string from .env file
	err := godotenv.Load()
	if err != nil {
		logger.Warn(fmt.Sprintf("failed to load env, %v", err))
	}

	conn, err := db.Connect()
	if err != nil {
		logger.Fatal("failed connecting to db", err)
	}
	defer conn.Close()
	db.RunMigrations(conn)

	mqConnection, err := queue.Connect()
	if err != nil {
		logger.Fatal("failed connecting to queue service", err)
	}
	defer mqConnection.Close()
	producer, err := queue.CreateProducer(mqConnection, "backend")
	if err != nil {
		logger.Fatal("failed creating producer", err)
	}
	defer producer.Close()
	messages, err := queue.CreateConsumer(mqConnection, "federation")
	if err != nil {
		logger.Fatal("failed creating consumer", err)
	}
	go func() {
		for message := range messages {
			logger.Debug(fmt.Sprintf(" > Received message: %s\n", message.Body))
		}
	}()
	s := http.NewServer(logger)
	s.RunServer()

	// NOTE: os.Exit(0) will not run defer statements
}
