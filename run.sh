#!/bin/bash

# Quick start script for Card Management API

echo "Starting Card Management API Server..."
echo "Server will be available at http://localhost:8080"
echo ""
echo "To test the API, first login:"
echo "curl -X POST http://localhost:8080/auth/login \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"userID\": \"testuser\", \"password\": \"password123\"}'"
echo ""

go run cmd/server/main.go
