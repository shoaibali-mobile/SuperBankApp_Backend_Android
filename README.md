# Bank App Card Management API

A comprehensive Go backend service for managing credit cards, debit cards, and virtual cards with full CRUD operations, transaction limits, autopay, and settings management.

## Features

- **Credit Card Management**: View, update limits, manage autopay, update PIN, request add-on cards
- **Debit Card Management**: View, update limits, update PIN
- **Virtual Card Management**: Create, view, update, delete, regenerate, manage spending limits and status
- **Card Settings**: Comprehensive settings management for notifications, security, limits, statements, and authentication
- **Transaction Limits**: Manage domestic and international transaction limits for all card types
- **Authentication**: Bearer token-based authentication

## Prerequisites

- Go 1.21 or higher
- macOS (tested on macOS)

## Installation

1. Clone or navigate to the project directory:
```bash
cd /Users/shoaibali/Documents/BankAppMicroservices
```

2. Install dependencies:
```bash
go mod download
```

## Running the Server

Start the server:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication

#### POST /auth/login
Login to get authentication token.

**Request:**
```json
{
  "userID": "testuser",
  "password": "password123"
}
```

**Response:**
```json
{
  "userID": "testuser",
  "fullName": "John Doe",
  "email": "john.doe@example.com",
  "token": "your-auth-token-here",
  "expiryDate": "2024-01-16T00:00:00Z",
  "requiresPIN": false,
  "requiresOTP": false
}
```

### Credit Cards

- `GET /api/cards/credit` - Get all credit cards
- `GET /api/cards/credit/{cardId}` - Get credit card details
- `PUT /api/cards/credit/{cardId}/limits` - Update card limits
- `POST /api/cards/credit/{cardId}/autopay` - Enable autopay
- `PUT /api/cards/credit/{cardId}/autopay` - Update autopay
- `DELETE /api/cards/credit/{cardId}/autopay` - Disable autopay
- `POST /api/cards/credit/{cardId}/pin` - Update PIN
- `POST /api/cards/credit/{cardId}/addon` - Request add-on card

### Debit Cards

- `GET /api/cards/debit` - Get all debit cards
- `GET /api/cards/debit/{cardId}` - Get debit card details
- `PUT /api/cards/debit/{cardId}/limits` - Update card limits
- `POST /api/cards/debit/{cardId}/pin` - Update PIN

### Virtual Cards

- `GET /api/cards/virtual` - Get all virtual cards
- `GET /api/cards/virtual/{cardId}` - Get virtual card details
- `POST /api/cards/virtual` - Create virtual card
- `PUT /api/cards/virtual/{cardId}` - Update virtual card
- `DELETE /api/cards/virtual/{cardId}` - Delete virtual card
- `PUT /api/cards/virtual/{cardId}/spending-limit` - Update spending limit
- `PUT /api/cards/virtual/{cardId}/status` - Update card status
- `POST /api/cards/virtual/{cardId}/regenerate` - Regenerate card number
- `GET /api/cards/virtual/{cardId}/transactions` - Get transactions

### Card Settings

- `GET /api/cards/settings` - Get all settings
- `PUT /api/cards/settings/default` - Update default cards
- `PUT /api/cards/settings/security` - Update security settings
- `PUT /api/cards/settings/global-limits` - Update global limits
- `PUT /api/cards/settings/notifications` - Update notification settings
- `PUT /api/cards/settings/statement` - Update statement settings
- `PUT /api/cards/settings/pin` - Update PIN settings
- `PUT /api/cards/settings/authentication` - Update authentication settings

### Transaction Limits

- `GET /api/cards/{cardId}/limits` - Get transaction limits
- `PUT /api/cards/{cardId}/limits/domestic` - Update domestic limits
- `PUT /api/cards/{cardId}/limits/international` - Update international limits

## Testing

All endpoints require authentication. First, login to get a token:

```bash
# Login
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "userID": "testuser",
    "password": "password123"
  }' | jq -r '.token')

echo "Token: $TOKEN"

# Get credit cards
curl -X GET http://localhost:8080/api/cards/credit \
  -H "Authorization: Bearer $TOKEN"
```

See `CARD_API_CURL_COMMANDS.md` for complete testing examples.

## Project Structure

```
BankAppMicroservices/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── handlers/            # HTTP handlers
│   │   ├── auth.go         # Authentication handler
│   │   ├── credit.go       # Credit card handlers
│   │   ├── debit.go        # Debit card handlers
│   │   ├── virtual.go      # Virtual card handlers
│   │   ├── settings.go     # Settings handlers
│   │   ├── limits.go       # Transaction limits handlers
│   │   └── common.go       # Common helper functions
│   ├── middleware/         # HTTP middleware
│   │   └── auth.go         # Authentication middleware
│   ├── models/             # Data models
│   │   └── models.go       # All struct definitions
│   └── store/              # In-memory data store
│       └── store.go        # Store implementation
├── go.mod                  # Go module file
└── README.md              # This file
```

## Default Test Data

The server initializes with default test data:
- **User**: `testuser` / `password123`
- **Credit Card**: One default credit card
- **Debit Card**: One default debit card
- **Virtual Card**: One default virtual card
- **Settings**: Default card settings

## Notes

- All data is stored in-memory and will be reset when the server restarts
- All endpoints require Bearer token authentication (except `/auth/login`)
- Card numbers and CVVs are masked in responses for security
- Virtual card status can be: "Active", "Frozen", or "Cancelled"
- Virtual card expiry periods: "3 Months", "6 Months", "12 Months" or custom date

## Development

To build the application:
```bash
go build -o bin/server cmd/server/main.go
```

To run the built binary:
```bash
./bin/server
```

## License

This is a sample project for demonstration purposes.
