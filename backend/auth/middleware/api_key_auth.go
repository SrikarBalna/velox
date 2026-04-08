package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/rishik92/velox/auth/service"
)

const (
	ScopesKey contextKey = "scopes"
)

type APIKeyAuthMiddleware struct {
	svc *service.APIKeyService
}

func NewAPIKeyAuthMiddleware(svc *service.APIKeyService) *APIKeyAuthMiddleware {
	return &APIKeyAuthMiddleware{svc: svc}
}

// Authenticate protects endpoints with an API key.
func (m *APIKeyAuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		plaintextKey := parts[1]

		// Validate key (CSPRNG hash, CRC32, etc. handled by service)
		apiKey, err := m.svc.ValidateKey(plaintextKey)
		if err != nil {
			http.Error(w, "invalid API key", http.StatusUnauthorized)
			return
		}

		// Inject userID and scopes into context
		ctx := context.WithValue(r.Context(), UserIDKey, apiKey.UserID)
		ctx = context.WithValue(ctx, ScopesKey, apiKey.Scopes)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CheckScope verifies if the API key has the required scope.
func CheckScope(requiredScope string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scopes, ok := r.Context().Value(ScopesKey).([]string)
		if !ok {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		hasScope := false
		for _, s := range scopes {
			if s == requiredScope {
				hasScope = true
				break
			}
		}

		if !hasScope {
			http.Error(w, "insufficient scopes", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
