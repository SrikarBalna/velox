package repository

import (
	"database/sql"

	"github.com/rishik92/velox/auth/model"
)

type APILogRepository struct {
	db *sql.DB
}

func NewAPILogRepository(db *sql.DB) *APILogRepository {
	return &APILogRepository{db: db}
}

// CreateLog saves a single log entry.
func (r *APILogRepository) CreateLog(log *model.APILog) error {
	query := `
		INSERT INTO api_logs (api_key_id, submission_id, endpoint, method, status_code, duration_ms, overall_state, language, error_msg)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(query, log.APIKeyID, log.SubmissionID, log.Endpoint, log.Method, log.StatusCode, log.DurationMs, log.OverallState, log.Language, log.ErrorMsg)
	return err
}

// UpdateLogResult updates the result of a submission log based on submission_id.
func (r *APILogRepository) UpdateLogResult(submissionID string, overallState string, errorMsg string) error {
	query := `
		UPDATE api_logs 
		SET overall_state = $1, error_msg = $2 
		WHERE submission_id = $3`
	_, err := r.db.Exec(query, overallState, errorMsg, submissionID)
	return err
}

// GetStats calculates metrics for a given API key.
func (r *APILogRepository) GetStats(apiKeyID string) (*model.APIKeyStats, error) {
	stats := &model.APIKeyStats{
		APIKeyID:    apiKeyID,
		ErrorCounts: make(map[string]int),
	}

	// 1. Total Requests
	err := r.db.QueryRow("SELECT COUNT(*) FROM api_logs WHERE api_key_id = $1", apiKeyID).Scan(&stats.TotalRequests)
	if err != nil {
		return nil, err
	}

	if stats.TotalRequests == 0 {
		return stats, nil
	}

	// 2. Peak RPM (Minute Bucketed)
	queryRPM := `
		SELECT MAX(count) FROM (
			SELECT COUNT(*) as count 
			FROM api_logs 
			WHERE api_key_id = $1 
			GROUP BY date_trunc('minute', created_at)
		) as subquery`
	err = r.db.QueryRow(queryRPM, apiKeyID).Scan(&stats.PeakRPM)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// 3. Peak RPD (Day Bucketed)
	queryRPD := `
		SELECT MAX(count) FROM (
			SELECT COUNT(*) as count 
			FROM api_logs 
			WHERE api_key_id = $1 
			GROUP BY date_trunc('day', created_at)
		) as subquery`
	err = r.db.QueryRow(queryRPD, apiKeyID).Scan(&stats.PeakRPD)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// 4. Success metrics and Error types
	queryErrors := `
		SELECT overall_state, COUNT(*) 
		FROM api_logs 
		WHERE api_key_id = $1 
		GROUP BY overall_state`
	rows, err := r.db.Query(queryErrors, apiKeyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var successCount int64
	for rows.Next() {
		var state string
		var count int
		if err := rows.Scan(&state, &count); err != nil {
			return nil, err
		}
		if state == "Accepted" {
			successCount = int64(count)
		} else if state != "" {
			stats.ErrorCounts[state] = count
		}
	}
	stats.SuccessRate = (float64(successCount) / float64(stats.TotalRequests)) * 100

	return stats, nil
}
