package model

import (
	"time"
)

// APILog represents a single API request log.
type APILog struct {
	ID           string    `json:"id"`
	APIKeyID     string    `json:"api_key_id"`
	SubmissionID string    `json:"submission_id,omitempty"`
	Endpoint     string    `json:"endpoint"`
	Method       string    `json:"method"`
	StatusCode   int       `json:"status_code"`
	DurationMs   int       `json:"duration_ms"`
	OverallState string    `json:"overall_state"` // e.g., "Accepted", "Compile Error", "Time Limit Exceeded"
	Language     string    `json:"language"`
	ErrorMsg     string    `json:"error_msg,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// APIKeyStats represents aggregated metrics for an API key.
type APIKeyStats struct {
	APIKeyID      string         `json:"api_key_id"`
	TotalRequests int64          `json:"total_requests"`
	PeakRPM       int64          `json:"peak_rpm"`
	PeakRPD       int64          `json:"peak_rpd"`
	SuccessRate   float64        `json:"success_rate"`
	ErrorCounts   map[string]int `json:"error_counts"`
}
