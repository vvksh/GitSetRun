package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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
	// Start the HTTP server
	http.HandleFunc("/webhook", handleWebhook)
	log.Printf("Server is running on http://localhost:%s/webhook", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
