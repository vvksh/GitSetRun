package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	checkInterval = 300 * time.Second // Interval to check for queued jobs
	ttlMinutes    = 60                // Time to live for the spot instance in minutes
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	// Get the server port
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port
	}
	ctx := context.Background()
	go CheckQueuedJobs(ctx, checkInterval)
	// Start the HTTP server
	http.HandleFunc("/webhook", handleWebhook)
	log.Printf("Server is running on http://localhost:%s/webhook", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
