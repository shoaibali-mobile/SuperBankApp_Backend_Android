package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	UserID       string    `json:"userID"`
	Password     string    `json:"-"`
	FullName     string    `json:"fullName"`
	Email        string    `json:"email"`
	Token        string    `json:"token,omitempty"`
	ExpiryDate   time.Time `json:"expiryDate,omitempty"`
	RequiresPIN  bool      `json:"requiresPIN"`
	RequiresOTP  bool      `json:"requiresOTP"`
}

// LoginRequest represents login request
type LoginRequest struct {
	UserID   string `json:"userID"`
	Password string `json:"password"`
}

// CreditCard represents a credit card
type CreditCard struct {
	ID                string  `json:"id"`
	CardNumber        string  `json:"cardNumber"`
	CVV               string  `json:"cvv"`
	ExpiryMonth       int     `json:"expiryMonth"`
	ExpiryYear        int     `json:"expiryYear"`
	CardholderName    string  `json:"cardholderName"`
	CardType          string  `json:"cardType"`
	RewardsPoints     int     `json:"rewardsPoints"`
	AvailableCredit   float64 `json:"availableCredit,omitempty"`
	TotalCredit       float64 `json:"totalCredit,omitempty"`
	OutstandingBalance float64 `json:"outstandingBalance,omitempty"`
	UserID            string  `json:"-"`
}

// DebitCard represents a debit card
type DebitCard struct {
	ID             string  `json:"id"`
	CardNumber     string  `json:"cardNumber"`
	CVV            string  `json:"cvv"`
	ExpiryMonth    int     `json:"expiryMonth"`
	ExpiryYear     int     `json:"expiryYear"`
	CardholderName string  `json:"cardholderName"`
	CardType       string  `json:"cardType"`
	AccountNumber  string  `json:"accountNumber"`
	BankName       string  `json:"bankName"`
	AccountBalance float64 `json:"accountBalance,omitempty"`
	UserID         string  `json:"-"`
}

// VirtualCard represents a virtual card
type VirtualCard struct {
	ID              string    `json:"id"`
	CardNumber      string    `json:"cardNumber"`
	CVV             string    `json:"cvv"`
	ExpiryMonth     int       `json:"expiryMonth"`
	ExpiryYear      int       `json:"expiryYear"`
	CardholderName  string    `json:"cardholderName"`
	CardType        string    `json:"cardType"`
	Nickname        string    `json:"nickname"`
	SpendingLimit   float64   `json:"spendingLimit"`
	RemainingBalance float64  `json:"remainingBalance"`
	CreatedAt       time.Time `json:"createdAt"`
	Status          string    `json:"status"`
	LinkedAccountID string    `json:"linkedAccountId"`
	UserID          string    `json:"-"`
}

// TransactionLimit represents a transaction limit
type TransactionLimit struct {
	ID          string  `json:"id,omitempty"`
	Type        string  `json:"type"`
	IsEnabled   bool    `json:"isEnabled"`
	CurrentLimit float64 `json:"currentLimit"`
	MaxLimit    float64 `json:"maxLimit,omitempty"`
	CanSetLimit bool    `json:"canSetLimit,omitempty"`
}

// LimitsRequest represents request to update limits
type LimitsRequest struct {
	DomesticLimits     []TransactionLimit `json:"domesticLimits"`
	InternationalLimits []TransactionLimit `json:"internationalLimits"`
}

// Autopay represents autopay settings
type Autopay struct {
	ID              string    `json:"autopayId,omitempty"`
	CardID          string    `json:"cardId"`
	AmountOption    string    `json:"amountOption"`
	LinkedAccountID string    `json:"linkedAccountId"`
	AutoPayEnabled  bool      `json:"autoPayEnabled,omitempty"`
	ActivationDate  time.Time `json:"activationDate,omitempty"`
	UserID          string    `json:"-"`
}

// AutopayRequest represents autopay request
type AutopayRequest struct {
	AmountOption    string `json:"amountOption"`
	LinkedAccountID string `json:"linkedAccountId"`
	AutoPayEnabled  bool   `json:"autoPayEnabled,omitempty"`
}

// PINUpdateRequest represents PIN update request
type PINUpdateRequest struct {
	NewPIN        string `json:"newPIN"`
	ConfirmPIN    string `json:"confirmPIN"`
	TermsAccepted bool   `json:"termsAccepted"`
}

// AddonCardRequest represents add-on card request
type AddonCardRequest struct {
	CustomerID   string `json:"customerID"`
	NameOnCard   string `json:"nameOnCard"`
	DateOfBirth  string `json:"dateOfBirth"`
	Relationship string `json:"relationship"`
}

// VirtualCardCreateRequest represents virtual card creation request
type VirtualCardCreateRequest struct {
	Nickname         string    `json:"nickname"`
	SpendingLimit    float64   `json:"spendingLimit"`
	CardType         string    `json:"cardType"`
	ExpiryPeriod     string    `json:"expiryPeriod"`
	CustomExpiryDate *time.Time `json:"customExpiryDate"`
	LinkedAccountID  string    `json:"linkedAccountId"`
}

// VirtualCardUpdateRequest represents virtual card update request
type VirtualCardUpdateRequest struct {
	Nickname *string `json:"nickname,omitempty"`
}

// SpendingLimitRequest represents spending limit update request
type SpendingLimitRequest struct {
	SpendingLimit float64 `json:"spendingLimit"`
}

// StatusRequest represents status update request
type StatusRequest struct {
	Status string `json:"status"`
}

// Transaction represents a card transaction
type Transaction struct {
	ID        string    `json:"id"`
	CardID    string    `json:"cardId"`
	Amount    float64   `json:"amount"`
	Merchant  string    `json:"merchant"`
	Date      time.Time `json:"date"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
}

// CardSettings represents card settings
type CardSettings struct {
	DefaultCreditCardID            string    `json:"defaultCreditCardId"`
	DefaultDebitCardID             string    `json:"defaultDebitCardId"`
	DefaultVirtualCardID           string    `json:"defaultVirtualCardId"`
	TransactionNotificationsEnabled bool     `json:"transactionNotificationsEnabled"`
	NotificationPreferences        []string `json:"notificationPreferences"`
	TransactionAmountThreshold     float64  `json:"transactionAmountThreshold"`
	InternationalTransactionAlerts bool     `json:"internationalTransactionAlerts"`
	ContactlessPaymentsEnabled     bool     `json:"contactlessPaymentsEnabled"`
	InternationalUsageEnabled      bool     `json:"internationalUsageEnabled"`
	OnlineTransactionsEnabled      bool     `json:"onlineTransactionsEnabled"`
	ATMWithdrawalsEnabled          bool     `json:"atmWithdrawalsEnabled"`
	DefaultDailyLimit              float64  `json:"defaultDailyLimit"`
	DefaultMonthlyLimit            float64  `json:"defaultMonthlyLimit"`
	StatementDelivery              string   `json:"statementDelivery"`
	StatementFrequency             string   `json:"statementFrequency"`
	EStatementEnabled              bool     `json:"eStatementEnabled"`
	BiometricAuthenticationEnabled bool     `json:"biometricAuthenticationEnabled"`
	TwoFactorAuthenticationEnabled bool     `json:"twoFactorAuthenticationEnabled"`
	TransactionAuthenticationRequired bool  `json:"transactionAuthenticationRequired"`
	PINForContactlessEnabled        bool    `json:"pinForContactlessEnabled"`
	UserID                          string  `json:"-"`
}

// DefaultCardsRequest represents default cards update request
type DefaultCardsRequest struct {
	DefaultCreditCardID  *string `json:"defaultCreditCardId,omitempty"`
	DefaultDebitCardID    *string `json:"defaultDebitCardId,omitempty"`
	DefaultVirtualCardID *string `json:"defaultVirtualCardId,omitempty"`
}

// SecuritySettingsRequest represents security settings update request
type SecuritySettingsRequest struct {
	ContactlessPaymentsEnabled *bool `json:"contactlessPaymentsEnabled,omitempty"`
	InternationalUsageEnabled  *bool `json:"internationalUsageEnabled,omitempty"`
	OnlineTransactionsEnabled  *bool `json:"onlineTransactionsEnabled,omitempty"`
	ATMWithdrawalsEnabled      *bool `json:"atmWithdrawalsEnabled,omitempty"`
}

// GlobalLimitsRequest represents global limits update request
type GlobalLimitsRequest struct {
	DefaultDailyLimit   *float64 `json:"defaultDailyLimit,omitempty"`
	DefaultMonthlyLimit *float64 `json:"defaultMonthlyLimit,omitempty"`
}

// NotificationSettingsRequest represents notification settings update request
type NotificationSettingsRequest struct {
	TransactionNotificationsEnabled *bool     `json:"transactionNotificationsEnabled,omitempty"`
	NotificationPreferences          *[]string `json:"notificationPreferences,omitempty"`
	TransactionAmountThreshold       *float64 `json:"transactionAmountThreshold,omitempty"`
	InternationalTransactionAlerts   *bool     `json:"internationalTransactionAlerts,omitempty"`
}

// StatementSettingsRequest represents statement settings update request
type StatementSettingsRequest struct {
	StatementDelivery  *string `json:"statementDelivery,omitempty"`
	StatementFrequency *string `json:"statementFrequency,omitempty"`
	EStatementEnabled  *bool   `json:"eStatementEnabled,omitempty"`
}

// PINSettingsRequest represents PIN settings update request
type PINSettingsRequest struct {
	PINForContactlessEnabled *bool `json:"pinForContactlessEnabled,omitempty"`
}

// AuthenticationSettingsRequest represents authentication settings update request
type AuthenticationSettingsRequest struct {
	BiometricAuthenticationEnabled    *bool `json:"biometricAuthenticationEnabled,omitempty"`
	TwoFactorAuthenticationEnabled    *bool `json:"twoFactorAuthenticationEnabled,omitempty"`
	TransactionAuthenticationRequired *bool `json:"transactionAuthenticationRequired,omitempty"`
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// CardsResponse represents cards list response
type CardsResponse struct {
	Cards []interface{} `json:"cards"`
}

// LimitsResponse represents limits response
type LimitsResponse struct {
	CardID              string            `json:"cardId"`
	DomesticLimits      []TransactionLimit `json:"domesticLimits,omitempty"`
	InternationalLimits []TransactionLimit `json:"internationalLimits,omitempty"`
}

// Pagination represents pagination info
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

// TransactionsResponse represents transactions response
type TransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
	Pagination   Pagination    `json:"pagination"`
}

// GenerateID generates a new UUID
func GenerateID() string {
	return uuid.New().String()
}
