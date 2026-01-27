package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"bankapp-microservices/internal/store"
)

type contextKey string

const UserIDKey contextKey = "userID"

// AuthMiddleware validates Bearer token
func AuthMiddleware(store *store.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondWithError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondWithError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			token := parts[1]
			user, exists := store.GetUserByToken(token)
			if !exists {
				respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, user.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"message": message,
	})
}
