package service

import (
	"log"

	"github.com/rishik92/velox/auth/model"
	"github.com/rishik92/velox/auth/repository"
)

type APILogService struct {
	repo    *repository.APILogRepository
	logChan chan *model.APILog
}

func NewAPILogService(repo *repository.APILogRepository) *APILogService {
	s := &APILogService{
		repo:    repo,
		logChan: make(chan *model.APILog, 1000), // Buffer for 1000 logs
	}
	go s.processLogs()
	return s
}

// Log pushes a log entry into the channel for async processing.
func (s *APILogService) Log(entry *model.APILog) {
	select {
	case s.logChan <- entry:
		// Log queued successfully
	default:
		// Channel full, drop log to avoid blocking the main thread
		log.Println("Warn: api_log channel full, dropping log")
	}
}

func (s *APILogService) processLogs() {
	for entry := range s.logChan {
		if err := s.repo.CreateLog(entry); err != nil {
			log.Printf("Error writing API log: %v", err)
		}
	}
}

// UpdateResult async updates a log entry with the final submission result.
func (s *APILogService) UpdateResult(submissionID string, overallState string, errorMsg string) {
	// For simplicity, we'll do this synchronously for now as it's an update.
	// In a real high-load system, this could also be queued.
	if err := s.repo.UpdateLogResult(submissionID, overallState, errorMsg); err != nil {
		log.Printf("Error updating API log result: %v", err)
	}
}

func (s *APILogService) GetStats(apiKeyID string) (*model.APIKeyStats, error) {
	return s.repo.GetStats(apiKeyID)
}

func (s *APILogService) Shutdown() {
	close(s.logChan)
}
