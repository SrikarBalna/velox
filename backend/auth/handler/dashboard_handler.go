package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rishik92/velox/auth/middleware"
	"github.com/rishik92/velox/auth/service"
)

type DashboardHandler struct {
	svc *service.DashboardService
}

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

// GetData returns user-specific dashboard data.
// @Summary Get Dashboard Data
// @Description Fetch profile and activity data for the authenticated user.
// @Tags Dashboard
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]any "Dashboard data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /dashboard [get]
func (h *DashboardHandler) GetData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract UserID from context (set by jwt_middleware.go)
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	data, err := h.svc.GetUserData(userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch dashboard data"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
