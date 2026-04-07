package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rishik92/velox/auth/db"
	"github.com/rishik92/velox/auth/handler"
	"github.com/rishik92/velox/auth/middleware"
	"github.com/rishik92/velox/auth/repository"
	"github.com/rishik92/velox/auth/service"
	"github.com/rishik92/velox/judge"
	veloxRedis "github.com/rishik92/velox/shared/redis"

	_ "github.com/rishik92/velox/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Velox API
// @version 1.0
// @description This is the Velox API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and then your personal token.

func main() {
	veloxRedis.Connect()

	// Connect to DB
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	fmt.Println("Connected to PostgreSQL database successfully.")

	// Run database migrations
	if err := db.RunMigrations(database); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}

	// Set up Auth Module
	repo := repository.NewUserRepository(database)
	svc := service.NewAuthService(repo)
	authHandler := handler.NewAuthHandler(svc)

	// Set up Dashboard Module
	dashboardSvc := service.NewDashboardService(repo)
	dashboardHandler := handler.NewDashboardHandler(dashboardSvc)

	// Auth Routes
	http.HandleFunc("/auth/signup", authHandler.Signup)
	http.HandleFunc("/auth/login", authHandler.Login)
	http.HandleFunc("/auth/logout", authHandler.Logout)

	// Protected Dashboard Route
	dashboardMux := http.NewServeMux()
	dashboardMux.HandleFunc("/dashboard", dashboardHandler.GetData)
	http.Handle("/dashboard", middleware.RequireAuth(dashboardMux))

	// General API
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/health", healthHandler)

	// Swagger documentation (only in development)
	env := os.Getenv("GO_ENV")
	if env == "" || env == "development" {
		http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
		fmt.Println("Swagger UI available at: http://localhost:8080/swagger/index.html")
	}

	fmt.Println("API Server running on :8080")
	handler := middleware.SecurityHeaders(corsMiddleware(http.DefaultServeMux))
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// submitHandler receives and queues a code submission.
// @Summary Submit Code
// @Description Queue a new code submission for judging.
// @Tags Judge
// @Accept json
// @Produce json
// @Param submission body judge.SubmissionRequest true "Submission Request"
// @Success 202 {string} string "Accepted: submission_id"
// @Failure 400 {string} string "Method not allowed, Invalid JSON, or Limits too high"
// @Failure 500 {string} string "Failed to queue submission"
// @Router /submit [post]
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

// statusHandler checks the status of a submission.
// @Summary Check Submission Status
// @Description Get the result or pending status of a submission by ID.
// @Tags Judge
// @Produce json
// @Param submission_id query string true "Submission ID"
// @Success 200 {object} judge.SubmissionResponse "Submission result or pending"
// @Failure 400 {string} string "Missing submission_id"
// @Router /status [get]
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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		env := os.Getenv("GO_ENV")

		var allowedOrigin string
		if env == "" || env == "development" {
			allowedOrigin = "*"
		} else {
			allowedOrigin = "https://example.com"
		}

		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}