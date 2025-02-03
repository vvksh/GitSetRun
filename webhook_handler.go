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
)

type GitHubWebhookPayload struct {
	Action     string `json:"action"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
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

	// Extract repository and job information
	repoFullName := payload.Repository.FullName
	jobID := payload.WorkflowJob.ID

	if payload.Action == "queued" && payload.WorkflowJob.Status == "queued" {
		log.Printf("Queued job detected: Repo=%s, JobId=%d queued", repoFullName, jobID)
		instanceID, err := provisionEC2(jobID, repoFullName)
		if err != nil {
			log.Printf("Failed to provision EC2 instance: %v", err)
			http.Error(w, "Failed to provision EC2 instance", http.StatusInternalServerError)
			return
		}
		log.Printf("EC2 instance provisioned: %s", instanceID)
	}

	// Process the webhook (customize as needed)
	log.Printf("Webhook received: Action=%s, WorkflowJobID=%d, Status=%s, URL=%s",
		payload.Action, payload.WorkflowJob.ID, payload.WorkflowJob.Status, payload.WorkflowJob.HTMLURL)

	// Respond to GitHub
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Webhook received")
}
