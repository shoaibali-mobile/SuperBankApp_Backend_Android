package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"bankapp-microservices/internal/middleware"
	"bankapp-microservices/internal/models"
	"bankapp-microservices/internal/store"
	"github.com/gorilla/mux"
)

type VirtualCardHandler struct {
	store *store.Store
}

func NewVirtualCardHandler(store *store.Store) *VirtualCardHandler {
	return &VirtualCardHandler{store: store}
}

func (h *VirtualCardHandler) GetVirtualCards(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	cards := h.store.GetVirtualCardsByUserID(userID)

	var maskedCards []interface{}
	for _, card := range cards {
		maskedCard := map[string]interface{}{
			"id":               card.ID,
			"cardNumber":       card.CardNumber,
			"cvv":              "***",
			"expiryMonth":      card.ExpiryMonth,
			"expiryYear":       card.ExpiryYear,
			"cardholderName":   card.CardholderName,
			"cardType":         card.CardType,
			"nickname":         card.Nickname,
			"spendingLimit":    card.SpendingLimit,
			"remainingBalance": card.RemainingBalance,
			"createdAt":        card.CreatedAt.Format(time.RFC3339),
			"status":           card.Status,
			"linkedAccountId":  card.LinkedAccountID,
		}
		maskedCards = append(maskedCards, maskedCard)
	}

	respondWithSuccess(w, models.CardsResponse{Cards: maskedCards})
}

func (h *VirtualCardHandler) GetVirtualCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	card, exists := h.store.GetVirtualCardByID(cardID)
	if !exists {
		respondWithError(w, http.StatusNotFound, "Virtual card not found")
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

func (h *VirtualCardHandler) CreateVirtualCard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req models.VirtualCardCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Calculate expiry date
	now := time.Now()
	var expiryDate time.Time
	if req.CustomExpiryDate != nil {
		expiryDate = *req.CustomExpiryDate
	} else {
		var months int
		switch req.ExpiryPeriod {
		case "3 Months":
			months = 3
		case "6 Months":
			months = 6
		case "12 Months":
			months = 12
		default:
			months = 3
		}
		expiryDate = now.AddDate(0, months, 0)
	}

	card := &models.VirtualCard{
		ID:               models.GenerateID(),
		CardNumber:       generateCardNumber(),
		CVV:              generateCVV(),
		ExpiryMonth:      int(expiryDate.Month()),
		ExpiryYear:       expiryDate.Year(),
		CardholderName:   "John Doe", // Get from user
		CardType:         req.CardType,
		Nickname:         req.Nickname,
		SpendingLimit:    req.SpendingLimit,
		RemainingBalance: req.SpendingLimit,
		CreatedAt:        now,
		Status:           "Active",
		LinkedAccountID:  req.LinkedAccountID,
		UserID:           userID,
	}

	h.store.CreateVirtualCard(card)

	respondWithSuccess(w, card, "Virtual card created successfully")
}

func (h *VirtualCardHandler) UpdateVirtualCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	card, exists := h.store.GetVirtualCardByID(cardID)
	if !exists {
		respondWithError(w, http.StatusNotFound, "Virtual card not found")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	if card.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req models.VirtualCardUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Nickname != nil {
		card.Nickname = *req.Nickname
	}

	h.store.UpdateVirtualCard(card)

	respondWithSuccess(w, card, "Virtual card updated successfully")
}

func (h *VirtualCardHandler) UpdateSpendingLimit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	card, exists := h.store.GetVirtualCardByID(cardID)
	if !exists {
		respondWithError(w, http.StatusNotFound, "Virtual card not found")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	if card.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req models.SpendingLimitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	card.SpendingLimit = req.SpendingLimit
	// Adjust remaining balance proportionally
	if card.SpendingLimit > 0 {
		card.RemainingBalance = (card.RemainingBalance / card.SpendingLimit) * req.SpendingLimit
	}

	h.store.UpdateVirtualCard(card)

	respondWithSuccess(w, map[string]interface{}{
		"cardId":           cardID,
		"spendingLimit":    card.SpendingLimit,
		"remainingBalance": card.RemainingBalance,
	}, "Spending limit updated successfully")
}

func (h *VirtualCardHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	card, exists := h.store.GetVirtualCardByID(cardID)
	if !exists {
		respondWithError(w, http.StatusNotFound, "Virtual card not found")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	if card.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req models.StatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	validStatuses := map[string]bool{"Active": true, "Frozen": true, "Cancelled": true}
	if !validStatuses[req.Status] {
		respondWithError(w, http.StatusBadRequest, "Invalid status. Must be Active, Frozen, or Cancelled")
		return
	}

	card.Status = req.Status
	h.store.UpdateVirtualCard(card)

	respondWithSuccess(w, map[string]interface{}{
		"cardId": cardID,
		"status": card.Status,
	}, "Card status updated successfully")
}

func (h *VirtualCardHandler) DeleteVirtualCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	card, exists := h.store.GetVirtualCardByID(cardID)
	if !exists {
		respondWithError(w, http.StatusNotFound, "Virtual card not found")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	if card.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	h.store.DeleteVirtualCard(cardID)

	respondWithSuccess(w, nil, "Virtual card deleted successfully")
}

func (h *VirtualCardHandler) RegenerateCard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	card, exists := h.store.GetVirtualCardByID(cardID)
	if !exists {
		respondWithError(w, http.StatusNotFound, "Virtual card not found")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	if card.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	card.CardNumber = generateCardNumber()
	card.CVV = generateCVV()
	h.store.UpdateVirtualCard(card)

	respondWithSuccess(w, map[string]interface{}{
		"cardId":        cardID,
		"newCardNumber": card.CardNumber,
		"newCVV":        card.CVV,
		"expiryMonth":   card.ExpiryMonth,
		"expiryYear":    card.ExpiryYear,
	}, "Card number regenerated successfully")
}

func (h *VirtualCardHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardId"]

	card, exists := h.store.GetVirtualCardByID(cardID)
	if !exists {
		respondWithError(w, http.StatusNotFound, "Virtual card not found")
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	if card.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Access denied")
		return
	}

	transactions := h.store.GetTransactionsByCardID(cardID)

	// Parse query parameters
	page := 1
	limit := 20
	if p := r.URL.Query().Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	// Filter by date if provided
	var filtered []*models.Transaction
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")

	for _, txn := range transactions {
		if startDateStr != "" {
			if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
				if txn.Date.Before(startDate) {
					continue
				}
			}
		}
		if endDateStr != "" {
			if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
				endDate = endDate.Add(24 * time.Hour)
				if txn.Date.After(endDate) {
					continue
				}
			}
		}
		filtered = append(filtered, txn)
	}

	// Paginate
	total := len(filtered)
	totalPages := (total + limit - 1) / limit
	start := (page - 1) * limit
	end := start + limit
	if end > total {
		end = total
	}

	var paginated []models.Transaction
	if start < total {
		for _, txn := range filtered[start:end] {
			paginated = append(paginated, *txn)
		}
	}

	respondWithSuccess(w, models.TransactionsResponse{
		Transactions: paginated,
		Pagination: models.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func generateCardNumber() string {
	return "4532" + generateRandomDigits(12)
}

func generateCVV() string {
	return generateRandomDigits(3)
}

func generateRandomDigits(n int) string {
	rand.Seed(time.Now().UnixNano())
	result := ""
	for i := 0; i < n; i++ {
		result += strconv.Itoa(rand.Intn(10))
	}
	return result
}
