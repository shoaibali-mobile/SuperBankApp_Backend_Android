package handlers

import (
	"encoding/json"
	"net/http"

	"bankapp-microservices/internal/middleware"
	"bankapp-microservices/internal/models"
	"bankapp-microservices/internal/store"
	"github.com/gorilla/mux"
)

type LimitsHandler struct {
	store *store.Store
}

func NewLimitsHandler(store *store.Store) *LimitsHandler {
	return &LimitsHandler{store: store}
}

func (h *LimitsHandler) GetLimits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	// Check if card exists (credit, debit, or virtual)
	creditCard, creditExists := h.store.GetCreditCardByID(cardID)
	debitCard, debitExists := h.store.GetDebitCardByID(cardID)
	virtualCard, virtualExists := h.store.GetVirtualCardByID(cardID)

	if !creditExists && !debitExists && !virtualExists {
		respondWithError(w, http.StatusNotFound, "Card not found")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	var cardUserID string
	if creditExists {
		cardUserID = creditCard.UserID
	} else if debitExists {
		cardUserID = debitCard.UserID
	} else {
		cardUserID = virtualCard.UserID
	}

	if cardUserID != userID {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	limits, exists := h.store.GetCardLimits(cardID)
	if !exists {
		// Return default limits structure
		limits = &models.LimitsRequest{
			DomesticLimits: []models.TransactionLimit{
				{Type: "ATM Cash Withdrawal", IsEnabled: true, CurrentLimit: 50000, MaxLimit: 100000, CanSetLimit: true},
				{Type: "Online", IsEnabled: true, CurrentLimit: 200000, MaxLimit: 500000, CanSetLimit: true},
				{Type: "Merchant Outlets (POS)", IsEnabled: true, CurrentLimit: 150000, MaxLimit: 300000, CanSetLimit: true},
				{Type: "Contactless", IsEnabled: true, CurrentLimit: 5000, MaxLimit: 10000, CanSetLimit: true},
			},
			InternationalLimits: []models.TransactionLimit{
				{Type: "ATM Cash Withdrawal", IsEnabled: false, CurrentLimit: 0, MaxLimit: 50000, CanSetLimit: true},
				{Type: "Online", IsEnabled: true, CurrentLimit: 100000, MaxLimit: 200000, CanSetLimit: true},
				{Type: "Merchant Outlets (POS)", IsEnabled: false, CurrentLimit: 0, MaxLimit: 100000, CanSetLimit: true},
				{Type: "Contactless", IsEnabled: false, CurrentLimit: 0, MaxLimit: 5000, CanSetLimit: true},
			},
		}
	}

	// Add IDs and max limits
	for i := range limits.DomesticLimits {
		if limits.DomesticLimits[i].ID == "" {
			limits.DomesticLimits[i].ID = models.GenerateID()
		}
		if limits.DomesticLimits[i].MaxLimit == 0 {
			limits.DomesticLimits[i].MaxLimit = limits.DomesticLimits[i].CurrentLimit
		}
		limits.DomesticLimits[i].CanSetLimit = true
	}
	for i := range limits.InternationalLimits {
		if limits.InternationalLimits[i].ID == "" {
			limits.InternationalLimits[i].ID = models.GenerateID()
		}
		if limits.InternationalLimits[i].MaxLimit == 0 {
			limits.InternationalLimits[i].MaxLimit = limits.InternationalLimits[i].CurrentLimit
		}
		limits.InternationalLimits[i].CanSetLimit = true
	}

	respondWithSuccess(w, models.LimitsResponse{
		CardID:              cardID,
		DomesticLimits:      limits.DomesticLimits,
		InternationalLimits: limits.InternationalLimits,
	})
}

func (h *LimitsHandler) UpdateDomesticLimits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	// Check if card exists
	creditCard, creditExists := h.store.GetCreditCardByID(cardID)
	debitCard, debitExists := h.store.GetDebitCardByID(cardID)
	virtualCard, virtualExists := h.store.GetVirtualCardByID(cardID)

	if !creditExists && !debitExists && !virtualExists {
		respondWithError(w, http.StatusNotFound, "Card not found")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	var cardUserID string
	if creditExists {
		cardUserID = creditCard.UserID
	} else if debitExists {
		cardUserID = debitCard.UserID
	} else {
		cardUserID = virtualCard.UserID
	}

	if cardUserID != userID {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req struct {
		Limits []models.TransactionLimit `json:"limits"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	limits, exists := h.store.GetCardLimits(cardID)
	if !exists {
		limits = &models.LimitsRequest{}
	}
	limits.DomesticLimits = req.Limits
	h.store.SetCardLimits(cardID, limits)

	respondWithSuccess(w, map[string]interface{}{
		"cardId":         cardID,
		"domesticLimits": limits.DomesticLimits,
	}, "Domestic limits updated successfully")
}

func (h *LimitsHandler) UpdateInternationalLimits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	// Check if card exists
	creditCard, creditExists := h.store.GetCreditCardByID(cardID)
	debitCard, debitExists := h.store.GetDebitCardByID(cardID)
	virtualCard, virtualExists := h.store.GetVirtualCardByID(cardID)

	if !creditExists && !debitExists && !virtualExists {
		respondWithError(w, http.StatusNotFound, "Card not found")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	var cardUserID string
	if creditExists {
		cardUserID = creditCard.UserID
	} else if debitExists {
		cardUserID = debitCard.UserID
	} else {
		cardUserID = virtualCard.UserID
	}

	if cardUserID != userID {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req struct {
		Limits []models.TransactionLimit `json:"limits"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	limits, exists := h.store.GetCardLimits(cardID)
	if !exists {
		limits = &models.LimitsRequest{}
	}
	limits.InternationalLimits = req.Limits
	h.store.SetCardLimits(cardID, limits)

	respondWithSuccess(w, map[string]interface{}{
		"cardId":              cardID,
		"internationalLimits": limits.InternationalLimits,
	}, "International limits updated successfully")
}
