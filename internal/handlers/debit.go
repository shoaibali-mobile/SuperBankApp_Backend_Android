package handlers

import (
	"encoding/json"
	"net/http"

	"bankapp-microservices/internal/middleware"
	"bankapp-microservices/internal/models"
	"bankapp-microservices/internal/store"
	"github.com/gorilla/mux"
)

type DebitCardHandler struct {
	store *store.Store
}

func NewDebitCardHandler(store *store.Store) *DebitCardHandler {
	return &DebitCardHandler{store: store}
}

func (h *DebitCardHandler) GetDebitCards(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	cards := h.store.GetDebitCardsByUserID(userID)

	var maskedCards []interface{}
	for _, card := range cards {
		maskedCard := map[string]interface{}{
			"id":             card.ID,
			"cardNumber":     card.CardNumber,
			"cvv":            "***",
			"expiryMonth":    card.ExpiryMonth,
			"expiryYear":     card.ExpiryYear,
			"cardholderName": card.CardholderName,
			"cardType":       card.CardType,
			"accountNumber":  card.AccountNumber,
			"bankName":       card.BankName,
		}
		maskedCards = append(maskedCards, maskedCard)
	}

	respondWithSuccess(w, models.CardsResponse{Cards: maskedCards})
}

func (h *DebitCardHandler) GetDebitCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	card, exists := h.store.GetDebitCardByID(cardID)
	if !exists {
		respondWithError(w, http.StatusNotFound, "Debit card not found")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	if card.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	card.CVV = "***"
	respondWithSuccess(w, card)
}

func (h *DebitCardHandler) UpdateLimits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	card, exists := h.store.GetDebitCardByID(cardID)
	if !exists {
		respondWithError(w, http.StatusNotFound, "Debit card not found")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	if card.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req models.LimitsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	h.store.SetCardLimits(cardID, &req)

	respondWithSuccess(w, map[string]interface{}{
		"cardId": cardID,
		"limits": req,
	}, "Debit card limits updated successfully")
}

func (h *DebitCardHandler) UpdatePIN(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	card, exists := h.store.GetDebitCardByID(cardID)
	if !exists {
		respondWithError(w, http.StatusNotFound, "Debit card not found")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	if card.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req models.PINUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.NewPIN != req.ConfirmPIN {
		respondWithError(w, http.StatusBadRequest, "PINs do not match")
		return
	}

	if !req.TermsAccepted {
		respondWithError(w, http.StatusBadRequest, "Terms must be accepted")
		return
	}

	respondWithSuccess(w, nil, "PIN updated successfully")
}
