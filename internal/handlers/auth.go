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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
