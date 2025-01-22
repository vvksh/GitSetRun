package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
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

type GitHubWebhookPayload struct {
	Action      string `json:"action"`
	WorkflowJob struct {
		ID      int64  `json:"id"`
		Status  string `json:"status"`
		HTMLURL string `json:"html_url"`
	} `json:"workflow_job"`
}

func verifySignature(secret, payload []byte, signature string) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	expectedMAC := mac.Sum(nil)
	expectedSignature := "sha256=" + hex.EncodeToString(expectedMAC)
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {

	// Read the webhook secret from environment variables
	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")

	if secret == "" {
		http.Error(w, "GITHUB_WEBHOOK_SECRET missing", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	fmt.Printf("body: %s\n", body)
	if err != nil {
		http.Error(w, "Failed to read request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		http.Error(w, "signature missing", http.StatusBadRequest)
		return
	}

	if !verifySignature([]byte(secret), body, signature) {
		http.Error(w, "signature mismatch", http.StatusForbidden)
		return
	}

	// Parse the JSON payload
	var payload GitHubWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	// Process the webhook (customize as needed)
	log.Printf("Webhook received: Action=%s, WorkflowJobID=%d, Status=%s, URL=%s",
		payload.Action, payload.WorkflowJob.ID, payload.WorkflowJob.Status, payload.WorkflowJob.HTMLURL)

	// Respond to GitHub
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Webhook received")
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
