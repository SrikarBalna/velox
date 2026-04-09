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
	svc    *service.APIKeyService
	logSvc *service.APILogService
}

func NewAPIKeyHandler(svc *service.APIKeyService, logSvc *service.APILogService) *APIKeyHandler {
	return &APIKeyHandler{svc: svc, logSvc: logSvc}
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

type UpdateKeyRequest struct {
	Name string `json:"name"`
}

// GenerateKey is an HTTP handler to create a new API key for the current user.
// @Summary Generate API Key
// @Description Create a new API key with specific scopes and optional expiration.
// @Tags APIKeys
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body GenerateKeyRequest true "API Key Generation Request"
// @Success 201 {object} GenerateKeyResponse "API key created"
// @Failure 400 {object} errorResponse "Invalid request"
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /auth/api-keys [post]
func (h *APIKeyHandler) GenerateKey(w http.ResponseWriter, r *http.Request) {
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

// ListKeys returns all API keys for the current user.
// @Summary List API Keys
// @Description Retrieve a list of all API keys created by the current user.
// @Tags APIKeys
// @Produce json
// @Security Bearer
// @Success 200 {array} model.APIKey "List of API keys"
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /auth/api-keys [get]
func (h *APIKeyHandler) ListKeys(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	keys, err := h.svc.ListKeys(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to list keys: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(keys)
}

// UpdateKey updates the name of an API key.
// @Summary Update API Key Name
// @Description Rename an existing API key.
// @Tags APIKeys
// @Accept json
// @Security Bearer
// @Param id query string true "API Key UUID"
// @Param request body UpdateKeyRequest true "Update Request"
// @Success 204 "Key updated"
// @Failure 400 {object} errorResponse "Missing ID or invalid JSON"
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /auth/api-keys [patch]
func (h *APIKeyHandler) UpdateKey(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	var req UpdateKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdateKeyName(id, userID, req.Name); err != nil {
		http.Error(w, fmt.Sprintf("failed to update key: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteKey deletes an API key.
// @Summary Delete API Key
// @Description Permanently revoke and delete an API key.
// @Tags APIKeys
// @Security Bearer
// @Param id query string true "API Key UUID"
// @Success 204 "Key deleted"
// @Failure 400 {object} errorResponse "Missing ID"
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /auth/api-keys [delete]
func (h *APIKeyHandler) DeleteKey(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteKey(id, userID); err != nil {
		http.Error(w, fmt.Sprintf("failed to delete key: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetStats returns usage statistics for an API key.
// @Summary Get API Key Stats
// @Description Retrieve performance metrics and usage peaks for a specific API key.
// @Tags APIKeys
// @Produce json
// @Security Bearer
// @Param id query string true "API Key ID"
// @Success 200 {object} model.APIKeyStats "API Key Statistics"
// @Failure 401 {object} errorResponse "Unauthorized"
// @Failure 404 {object} errorResponse "Not Found"
// @Router /auth/api-keys/stats [get]
func (h *APIKeyHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	stats, err := h.logSvc.GetStats(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get stats: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}
