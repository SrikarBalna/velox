package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rishik92/velox/auth/middleware"
	"github.com/rishik92/velox/auth/service"
)

type APIKeyHandler struct {
	svc *service.APIKeyService
}

func NewAPIKeyHandler(svc *service.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{svc: svc}
}

type GenerateKeyRequest struct {
	Name      string     `json:"name"`
	Scopes    []string   `json:"scopes"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type GenerateKeyResponse struct {
	Key         string     `json:"key"`
	ID          string     `json:"id"`
	DisplayHint string     `json:"display_hint"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

// GenerateKey is an HTTP handler to create a new API key for the current user.
func (h *APIKeyHandler) GenerateKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get userID from context (set by RequireAuth middleware)
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req GenerateKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	// Default scopes if none provided
	if len(req.Scopes) == 0 {
		req.Scopes = []string{"submit", "status"}
	}

	fullKey, apiKey, err := h.svc.GenerateKey(userID, req.Name, req.Scopes, req.ExpiresAt)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to generate key: %v", err), http.StatusInternalServerError)
		return
	}

	resp := GenerateKeyResponse{
		Key:         fullKey,
		ID:          apiKey.ID,
		DisplayHint: apiKey.DisplayHint,
		CreatedAt:   apiKey.CreatedAt,
		ExpiresAt:   apiKey.ExpiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
