package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rishik92/velox/judge"
	veloxRedis "github.com/rishik92/velox/shared/redis"
)

func main() {
	veloxRedis.Connect()

	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/status", statusHandler)

	fmt.Println("API Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req judge.SubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.TimeLimitMs > 5000 || req.MemoryLimitKb > 512000 {
		http.Error(w, "Invalid request: Limits too high (Max 5s, 512MB)", http.StatusBadRequest)
		return
	}

	// Generate unique ID
	req.SubmissionID = uuid.New().String()

	// Push to Redis
	raw, _ := json.Marshal(req)
	if err := veloxRedis.PushResult("submissions", string(raw)); err != nil {
		http.Error(w, "Failed to queue submission", http.StatusInternalServerError)
		return
	}

	// Respond immediately with ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_, _ = fmt.Fprintf(w, `{"submission_id": "%s"}`, req.SubmissionID)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	subID := r.URL.Query().Get("submission_id")
	if subID == "" {
		http.Error(w, "Missing submission_id", http.StatusBadRequest)
		return
	}

	// Check Redis for result
	resultQueue := "results:" + subID
	raw, found := veloxRedis.PopSubmission(resultQueue, 1*time.Second)

	if !found {
		// Still processing or not found
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"status": "pending"}`)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(raw))
}