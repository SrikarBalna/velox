package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rishik92/velox/auth/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type signupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type successResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// Signup handles user registration.
// @Summary User Signup
// @Description Register a new user with name, email, and password.
// @Tags Auth
// @Accept json
// @Produce json
// @Param signup body signupRequest true "Signup Request"
// @Success 201 {object} successResponse "User created"
// @Failure 400 {object} errorResponse "Invalid request or validation error"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /auth/signup [post]
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.svc.Signup(req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrNameRequired) || errors.Is(err, service.ErrEmailTaken) || errors.Is(err, service.ErrInvalidEmail) || errors.Is(err, service.ErrPasswordTooShort) {
			h.respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		// Log actual error internally
		h.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.respondJSON(w, http.StatusCreated, successResponse{
		Message: "user created successfully",
		Data: map[string]string{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

// Login handles user authentication.
// @Summary User Login
// @Description Authenticate a user and return a JWT token.
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body loginRequest true "Login Request"
// @Success 200 {object} successResponse "Login successful"
// @Failure 401 {object} errorResponse "Invalid credentials"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) || errors.Is(err, service.ErrInvalidEmail) || errors.Is(err, service.ErrPasswordTooShort) {
			h.respondError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		// Log actual error internally
		h.respondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.respondJSON(w, http.StatusOK, successResponse{
		Message: "login successful",
		Data: map[string]string{
			"token": token,
		},
	})
}

// Logout handles user logout.
// @Summary User Logout
// @Description Invalidate the session (stateless JWT logout).
// @Tags Auth
// @Produce json
// @Success 200 {object} successResponse "Logout successful"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// In a stateless JWT setup, logout operations are generally handled client-side
	// by deleting the token. For blacklisting, we would need a storage mechanism 
	// like Redis to check tokens against a blacklist during validation.
	h.respondJSON(w, http.StatusOK, successResponse{
		Message: "logout successful (please discard your token locally)",
	})
}

func (h *AuthHandler) respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (h *AuthHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, errorResponse{Error: message})
}
