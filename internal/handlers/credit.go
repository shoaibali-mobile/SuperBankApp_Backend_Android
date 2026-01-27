package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"bankapp-microservices/internal/middleware"
	"bankapp-microservices/internal/models"
	"bankapp-microservices/internal/store"
	"github.com/gorilla/mux"
)

type CreditCardHandler struct {
	store *store.Store
}

func NewCreditCardHandler(store *store.Store) *CreditCardHandler {
	return &CreditCardHandler{store: store}
}

func (h *CreditCardHandler) GetCreditCards(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	cards := h.store.GetCreditCardsByUserID(userID)

	// Mask CVV and card number
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
			"rewardsPoints":  card.RewardsPoints,
		}
		maskedCards = append(maskedCards, maskedCard)
	}

	respondWithSuccess(w, models.CardsResponse{Cards: maskedCards})
}

func (h *CreditCardHandler) GetCreditCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	card, exists := h.store.GetCreditCardByID(cardID)
	if !exists {
		respondWithError(w, http.StatusNotFound, "Credit card not found")
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

func (h *CreditCardHandler) UpdateLimits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	card, exists := h.store.GetCreditCardByID(cardID)
	if !exists {
		respondWithError(w, http.StatusNotFound, "Credit card not found")
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
	}, "Card limits updated successfully")
}

func (h *CreditCardHandler) EnableAutopay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	card, exists := h.store.GetCreditCardByID(cardID)
	if !exists {
		respondWithError(w, http.StatusNotFound, "Credit card not found")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	if card.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req models.AutopayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	autopay := &models.Autopay{
		ID:              models.GenerateID(),
		CardID:          cardID,
		AmountOption:    req.AmountOption,
		LinkedAccountID: req.LinkedAccountID,
		AutoPayEnabled:  req.AutoPayEnabled,
		ActivationDate:  time.Now(),
		UserID:          userID,
	}

	h.store.SetAutopay(autopay)

	respondWithSuccess(w, autopay, "Autopay enabled successfully")
}

func (h *CreditCardHandler) UpdateAutopay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	card, exists := h.store.GetCreditCardByID(cardID)
	if !exists {
		respondWithError(w, http.StatusNotFound, "Credit card not found")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	if card.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	autopay, exists := h.store.GetAutopayByCardID(cardID)
	if !exists {
		respondWithError(w, http.StatusNotFound, "Autopay not found")
		return
	}

	var req models.AutopayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	autopay.AmountOption = req.AmountOption
	autopay.LinkedAccountID = req.LinkedAccountID
	if req.AutoPayEnabled {
		autopay.AutoPayEnabled = req.AutoPayEnabled
	}

	h.store.SetAutopay(autopay)

	respondWithSuccess(w, nil, "Autopay settings updated successfully")
}

func (h *CreditCardHandler) DisableAutopay(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	card, exists := h.store.GetCreditCardByID(cardID)
	if !exists {
		respondWithError(w, http.StatusNotFound, "Credit card not found")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	if card.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	h.store.DeleteAutopay(cardID)

	respondWithSuccess(w, nil, "Autopay disabled successfully")
}

func (h *CreditCardHandler) UpdatePIN(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	card, exists := h.store.GetCreditCardByID(cardID)
	if !exists {
		respondWithError(w, http.StatusNotFound, "Credit card not found")
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

func (h *CreditCardHandler) RequestAddonCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	card, exists := h.store.GetCreditCardByID(cardID)
	if !exists {
		respondWithError(w, http.StatusNotFound, "Credit card not found")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	if card.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req models.AddonCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Calculate estimated delivery date (2 weeks from now)
	estimatedDelivery := time.Now().Add(14 * 24 * time.Hour)

	respondWithSuccess(w, map[string]interface{}{
		"requestId":            models.GenerateID(),
		"estimatedDeliveryDate": estimatedDelivery.Format(time.RFC3339),
	}, "Add-on card request submitted successfully")
}
