package handlers

import (
	"encoding/json"
	"net/http"

	"bankapp-microservices/internal/middleware"
	"bankapp-microservices/internal/models"
	"bankapp-microservices/internal/store"
)

type SettingsHandler struct {
	store *store.Store
}

func NewSettingsHandler(store *store.Store) *SettingsHandler {
	return &SettingsHandler{store: store}
}

func (h *SettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	settings, exists := h.store.GetCardSettings(userID)
	if !exists {
		// Create default settings
		settings = &models.CardSettings{
			UserID: userID,
		}
		h.store.UpdateCardSettings(settings)
	}

	respondWithSuccess(w, settings)
}

func (h *SettingsHandler) UpdateDefaultCards(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	settings, exists := h.store.GetCardSettings(userID)
	if !exists {
		settings = &models.CardSettings{UserID: userID}
	}

	var req models.DefaultCardsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.DefaultCreditCardID != nil {
		settings.DefaultCreditCardID = *req.DefaultCreditCardID
	}
	if req.DefaultDebitCardID != nil {
		settings.DefaultDebitCardID = *req.DefaultDebitCardID
	}
	if req.DefaultVirtualCardID != nil {
		settings.DefaultVirtualCardID = *req.DefaultVirtualCardID
	}

	h.store.UpdateCardSettings(settings)
	respondWithSuccess(w, nil, "Default cards updated successfully")
}

func (h *SettingsHandler) UpdateSecuritySettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	settings, exists := h.store.GetCardSettings(userID)
	if !exists {
		settings = &models.CardSettings{UserID: userID}
	}

	var req models.SecuritySettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.ContactlessPaymentsEnabled != nil {
		settings.ContactlessPaymentsEnabled = *req.ContactlessPaymentsEnabled
	}
	if req.InternationalUsageEnabled != nil {
		settings.InternationalUsageEnabled = *req.InternationalUsageEnabled
	}
	if req.OnlineTransactionsEnabled != nil {
		settings.OnlineTransactionsEnabled = *req.OnlineTransactionsEnabled
	}
	if req.ATMWithdrawalsEnabled != nil {
		settings.ATMWithdrawalsEnabled = *req.ATMWithdrawalsEnabled
	}

	h.store.UpdateCardSettings(settings)
	respondWithSuccess(w, nil, "Security settings updated successfully")
}

func (h *SettingsHandler) UpdateGlobalLimits(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	settings, exists := h.store.GetCardSettings(userID)
	if !exists {
		settings = &models.CardSettings{UserID: userID}
	}

	var req models.GlobalLimitsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.DefaultDailyLimit != nil {
		settings.DefaultDailyLimit = *req.DefaultDailyLimit
	}
	if req.DefaultMonthlyLimit != nil {
		settings.DefaultMonthlyLimit = *req.DefaultMonthlyLimit
	}

	h.store.UpdateCardSettings(settings)
	respondWithSuccess(w, nil, "Global transaction limits updated successfully")
}

func (h *SettingsHandler) UpdateNotificationSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	settings, exists := h.store.GetCardSettings(userID)
	if !exists {
		settings = &models.CardSettings{UserID: userID}
	}

	var req models.NotificationSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.TransactionNotificationsEnabled != nil {
		settings.TransactionNotificationsEnabled = *req.TransactionNotificationsEnabled
	}
	if req.NotificationPreferences != nil {
		settings.NotificationPreferences = *req.NotificationPreferences
	}
	if req.TransactionAmountThreshold != nil {
		settings.TransactionAmountThreshold = *req.TransactionAmountThreshold
	}
	if req.InternationalTransactionAlerts != nil {
		settings.InternationalTransactionAlerts = *req.InternationalTransactionAlerts
	}

	h.store.UpdateCardSettings(settings)
	respondWithSuccess(w, nil, "Notification preferences updated successfully")
}

func (h *SettingsHandler) UpdateStatementSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	settings, exists := h.store.GetCardSettings(userID)
	if !exists {
		settings = &models.CardSettings{UserID: userID}
	}

	var req models.StatementSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.StatementDelivery != nil {
		settings.StatementDelivery = *req.StatementDelivery
	}
	if req.StatementFrequency != nil {
		settings.StatementFrequency = *req.StatementFrequency
	}
	if req.EStatementEnabled != nil {
		settings.EStatementEnabled = *req.EStatementEnabled
	}

	h.store.UpdateCardSettings(settings)
	respondWithSuccess(w, nil, "Statement preferences updated successfully")
}

func (h *SettingsHandler) UpdatePINSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	settings, exists := h.store.GetCardSettings(userID)
	if !exists {
		settings = &models.CardSettings{UserID: userID}
	}

	var req models.PINSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.PINForContactlessEnabled != nil {
		settings.PINForContactlessEnabled = *req.PINForContactlessEnabled
	}

	h.store.UpdateCardSettings(settings)
	respondWithSuccess(w, nil, "PIN preferences updated successfully")
}

func (h *SettingsHandler) UpdateAuthenticationSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)
	settings, exists := h.store.GetCardSettings(userID)
	if !exists {
		settings = &models.CardSettings{UserID: userID}
	}

	var req models.AuthenticationSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.BiometricAuthenticationEnabled != nil {
		settings.BiometricAuthenticationEnabled = *req.BiometricAuthenticationEnabled
	}
	if req.TwoFactorAuthenticationEnabled != nil {
		settings.TwoFactorAuthenticationEnabled = *req.TwoFactorAuthenticationEnabled
	}
	if req.TransactionAuthenticationRequired != nil {
		settings.TransactionAuthenticationRequired = *req.TransactionAuthenticationRequired
	}

	h.store.UpdateCardSettings(settings)
	respondWithSuccess(w, nil, "Authentication settings updated successfully")
}
