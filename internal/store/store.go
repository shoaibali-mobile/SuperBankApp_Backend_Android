package store

import (
	"sync"
	"time"

	"bankapp-microservices/internal/models"
)

// Store represents in-memory data store
type Store struct {
	mu                sync.RWMutex
	users             map[string]*models.User
	tokens            map[string]string // token -> userID
	creditCards       map[string]*models.CreditCard
	debitCards        map[string]*models.DebitCard
	virtualCards      map[string]*models.VirtualCard
	autopays          map[string]*models.Autopay // cardID -> autopay
	cardLimits        map[string]*models.LimitsRequest // cardID -> limits
	cardSettings      map[string]*models.CardSettings // userID -> settings
	transactions      map[string][]*models.Transaction // cardID -> transactions
}

// NewStore creates a new store instance
func NewStore() *Store {
	store := &Store{
		users:        make(map[string]*models.User),
		tokens:       make(map[string]string),
		creditCards:  make(map[string]*models.CreditCard),
		debitCards:   make(map[string]*models.DebitCard),
		virtualCards: make(map[string]*models.VirtualCard),
		autopays:     make(map[string]*models.Autopay),
		cardLimits:   make(map[string]*models.LimitsRequest),
		cardSettings: make(map[string]*models.CardSettings),
		transactions: make(map[string][]*models.Transaction),
	}
	store.initDefaultData()
	return store
}

// initDefaultData initializes default test data
func (s *Store) initDefaultData() {
	// Create default user
	user := &models.User{
		UserID:      "testuser",
		Password:    "password123",
		FullName:    "John Doe",
		Email:       "john.doe@example.com",
		RequiresPIN: false,
		RequiresOTP: false,
	}
	s.users[user.UserID] = user

	// Create default credit card
	creditCard := &models.CreditCard{
		ID:                models.GenerateID(),
		CardNumber:        "4532123456789012",
		CVV:               "***",
		ExpiryMonth:       12,
		ExpiryYear:        2026,
		CardholderName:    "John Doe",
		CardType:          "Rupay",
		RewardsPoints:     1250,
		AvailableCredit:   50000.0,
		TotalCredit:       100000.0,
		OutstandingBalance: 50000.0,
		UserID:            user.UserID,
	}
	s.creditCards[creditCard.ID] = creditCard

	// Create default debit card
	debitCard := &models.DebitCard{
		ID:             models.GenerateID(),
		CardNumber:     "6529251234567890",
		CVV:            "***",
		ExpiryMonth:    10,
		ExpiryYear:     2028,
		CardholderName: "John Doe",
		CardType:       "Rupay",
		AccountNumber:  "50123456789012",
		BankName:       "HDFC Bank",
		AccountBalance: 50000.0,
		UserID:         user.UserID,
	}
	s.debitCards[debitCard.ID] = debitCard

	// Create default virtual card
	virtualCard := &models.VirtualCard{
		ID:               models.GenerateID(),
		CardNumber:       "4532123456789012",
		CVV:              "***",
		ExpiryMonth:      3,
		ExpiryYear:       2025,
		CardholderName:   "John Doe",
		CardType:         "Visa",
		Nickname:         "Netflix Subscription",
		SpendingLimit:    5000.0,
		RemainingBalance: 3200.0,
		CreatedAt:        time.Now(),
		Status:           "Active",
		LinkedAccountID:  "account-uuid",
		UserID:           user.UserID,
	}
	s.virtualCards[virtualCard.ID] = virtualCard

	// Create default card settings
	settings := &models.CardSettings{
		DefaultCreditCardID:              creditCard.ID,
		DefaultDebitCardID:               debitCard.ID,
		DefaultVirtualCardID:             virtualCard.ID,
		TransactionNotificationsEnabled: true,
		NotificationPreferences:         []string{"Push Notification", "Email"},
		TransactionAmountThreshold:      1000.0,
		InternationalTransactionAlerts:   true,
		ContactlessPaymentsEnabled:       true,
		InternationalUsageEnabled:        true,
		OnlineTransactionsEnabled:        true,
		ATMWithdrawalsEnabled:            true,
		DefaultDailyLimit:                50000.0,
		DefaultMonthlyLimit:              200000.0,
		StatementDelivery:                "Email",
		StatementFrequency:               "Monthly",
		EStatementEnabled:                true,
		BiometricAuthenticationEnabled:    true,
		TwoFactorAuthenticationEnabled:    false,
		TransactionAuthenticationRequired: true,
		PINForContactlessEnabled:          false,
		UserID:                           user.UserID,
	}
	s.cardSettings[user.UserID] = settings
}

// GetUserByID gets user by ID
func (s *Store) GetUserByID(userID string) (*models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, exists := s.users[userID]
	return user, exists
}

// GetUserByToken gets user by token
func (s *Store) GetUserByToken(token string) (*models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	userID, exists := s.tokens[token]
	if !exists {
		return nil, false
	}
	user, exists := s.users[userID]
	return user, exists
}

// SetToken sets token for user
func (s *Store) SetToken(token, userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[token] = userID
}

// GetCreditCardsByUserID gets all credit cards for a user
func (s *Store) GetCreditCardsByUserID(userID string) []*models.CreditCard {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var cards []*models.CreditCard
	for _, card := range s.creditCards {
		if card.UserID == userID {
			cards = append(cards, card)
		}
	}
	return cards
}

// GetCreditCardByID gets credit card by ID
func (s *Store) GetCreditCardByID(cardID string) (*models.CreditCard, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	card, exists := s.creditCards[cardID]
	return card, exists
}

// UpdateCreditCard updates credit card
func (s *Store) UpdateCreditCard(card *models.CreditCard) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.creditCards[card.ID] = card
}

// GetDebitCardsByUserID gets all debit cards for a user
func (s *Store) GetDebitCardsByUserID(userID string) []*models.DebitCard {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var cards []*models.DebitCard
	for _, card := range s.debitCards {
		if card.UserID == userID {
			cards = append(cards, card)
		}
	}
	return cards
}

// GetDebitCardByID gets debit card by ID
func (s *Store) GetDebitCardByID(cardID string) (*models.DebitCard, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	card, exists := s.debitCards[cardID]
	return card, exists
}

// UpdateDebitCard updates debit card
func (s *Store) UpdateDebitCard(card *models.DebitCard) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.debitCards[card.ID] = card
}

// GetVirtualCardsByUserID gets all virtual cards for a user
func (s *Store) GetVirtualCardsByUserID(userID string) []*models.VirtualCard {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var cards []*models.VirtualCard
	for _, card := range s.virtualCards {
		if card.UserID == userID {
			cards = append(cards, card)
		}
	}
	return cards
}

// GetVirtualCardByID gets virtual card by ID
func (s *Store) GetVirtualCardByID(cardID string) (*models.VirtualCard, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	card, exists := s.virtualCards[cardID]
	return card, exists
}

// CreateVirtualCard creates a new virtual card
func (s *Store) CreateVirtualCard(card *models.VirtualCard) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.virtualCards[card.ID] = card
}

// UpdateVirtualCard updates virtual card
func (s *Store) UpdateVirtualCard(card *models.VirtualCard) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.virtualCards[card.ID] = card
}

// DeleteVirtualCard deletes virtual card
func (s *Store) DeleteVirtualCard(cardID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.virtualCards, cardID)
}

// GetAutopayByCardID gets autopay by card ID
func (s *Store) GetAutopayByCardID(cardID string) (*models.Autopay, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	autopay, exists := s.autopays[cardID]
	return autopay, exists
}

// SetAutopay sets autopay for a card
func (s *Store) SetAutopay(autopay *models.Autopay) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.autopays[autopay.CardID] = autopay
}

// DeleteAutopay deletes autopay for a card
func (s *Store) DeleteAutopay(cardID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.autopays, cardID)
}

// GetCardLimits gets card limits
func (s *Store) GetCardLimits(cardID string) (*models.LimitsRequest, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	limits, exists := s.cardLimits[cardID]
	return limits, exists
}

// SetCardLimits sets card limits
func (s *Store) SetCardLimits(cardID string, limits *models.LimitsRequest) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cardLimits[cardID] = limits
}

// GetCardSettings gets card settings for user
func (s *Store) GetCardSettings(userID string) (*models.CardSettings, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	settings, exists := s.cardSettings[userID]
	return settings, exists
}

// UpdateCardSettings updates card settings
func (s *Store) UpdateCardSettings(settings *models.CardSettings) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cardSettings[settings.UserID] = settings
}

// GetTransactionsByCardID gets transactions for a card
func (s *Store) GetTransactionsByCardID(cardID string) []*models.Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	transactions, exists := s.transactions[cardID]
	if !exists {
		return []*models.Transaction{}
	}
	return transactions
}

// AddTransaction adds a transaction
func (s *Store) AddTransaction(transaction *models.Transaction) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.transactions[transaction.CardID] = append(s.transactions[transaction.CardID], transaction)
}
