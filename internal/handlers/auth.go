package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"bankapp-microservices/internal/models"
	"bankapp-microservices/internal/store"

	"github.com/google/uuid"
)

type AuthHandler struct {
	store *store.Store
}

func NewAuthHandler(store *store.Store) *AuthHandler {
	return &AuthHandler{store: store}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, exists := h.store.GetUserByID(req.UserID)
	if !exists || user.Password != req.Password {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate token
	token := uuid.New().String()
	user.Token = token
	user.ExpiryDate = time.Now().Add(24 * time.Hour)
	h.store.SetToken(token, user.UserID)

	// Format response to match Android expectations: { "user": { ... }, "tokens": { ... } }
	loginResponse := map[string]interface{}{
		"user": map[string]interface{}{
			"id":              user.UserID,
			"email":           user.Email,
			"name":            user.FullName,
			"createdAt":       time.Now().UnixMilli(),
			"accountStatus":   "ACTIVE",
			"isEmailVerified": true,
			"isPhoneVerified": true,
		},
		"tokens": map[string]interface{}{
			"accessToken":  token,
			"refreshToken": "dummy_refresh_token",
			"tokenType":    "Bearer",
			"expiresIn":    86400, // 24 hours in seconds
			"issuedAt":     time.Now().UnixMilli(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(loginResponse); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to encode response")
	}
}
