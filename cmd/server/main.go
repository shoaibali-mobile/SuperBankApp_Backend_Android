package main

import (
	"fmt"
	"log"
	"net/http"

	"bankapp-microservices/internal/handlers"
	"bankapp-microservices/internal/middleware"
	"bankapp-microservices/internal/store"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize store
	store := store.NewStore()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(store)
	creditHandler := handlers.NewCreditCardHandler(store)
	debitHandler := handlers.NewDebitCardHandler(store)
	virtualHandler := handlers.NewVirtualCardHandler(store)
	settingsHandler := handlers.NewSettingsHandler(store)
	limitsHandler := handlers.NewLimitsHandler(store)

	// Setup router
	r := mux.NewRouter()

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if req.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, req)
		})
	})

	// Public routes
	r.HandleFunc("/auth/login", authHandler.Login).Methods("POST", "OPTIONS")

	// Protected routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware(store))

	// Credit card routes
	creditRouter := api.PathPrefix("/cards/credit").Subrouter()
	creditRouter.HandleFunc("", creditHandler.GetCreditCards).Methods("GET")
	creditRouter.HandleFunc("/{cardId}", creditHandler.GetCreditCard).Methods("GET")
	creditRouter.HandleFunc("/{cardId}/limits", creditHandler.UpdateLimits).Methods("PUT")
	creditRouter.HandleFunc("/{cardId}/autopay", creditHandler.EnableAutopay).Methods("POST")
	creditRouter.HandleFunc("/{cardId}/autopay", creditHandler.UpdateAutopay).Methods("PUT")
	creditRouter.HandleFunc("/{cardId}/autopay", creditHandler.DisableAutopay).Methods("DELETE")
	creditRouter.HandleFunc("/{cardId}/pin", creditHandler.UpdatePIN).Methods("POST")
	creditRouter.HandleFunc("/{cardId}/addon", creditHandler.RequestAddonCard).Methods("POST")

	// Debit card routes
	debitRouter := api.PathPrefix("/cards/debit").Subrouter()
	debitRouter.HandleFunc("", debitHandler.GetDebitCards).Methods("GET")
	debitRouter.HandleFunc("/{cardId}", debitHandler.GetDebitCard).Methods("GET")
	debitRouter.HandleFunc("/{cardId}/limits", debitHandler.UpdateLimits).Methods("PUT")
	debitRouter.HandleFunc("/{cardId}/pin", debitHandler.UpdatePIN).Methods("POST")

	// Virtual card routes
	virtualRouter := api.PathPrefix("/cards/virtual").Subrouter()
	virtualRouter.HandleFunc("", virtualHandler.GetVirtualCards).Methods("GET")
	virtualRouter.HandleFunc("", virtualHandler.CreateVirtualCard).Methods("POST")
	virtualRouter.HandleFunc("/{cardId}", virtualHandler.GetVirtualCard).Methods("GET")
	virtualRouter.HandleFunc("/{cardId}", virtualHandler.UpdateVirtualCard).Methods("PUT")
	virtualRouter.HandleFunc("/{cardId}", virtualHandler.DeleteVirtualCard).Methods("DELETE")
	virtualRouter.HandleFunc("/{cardId}/spending-limit", virtualHandler.UpdateSpendingLimit).Methods("PUT")
	virtualRouter.HandleFunc("/{cardId}/status", virtualHandler.UpdateStatus).Methods("PUT")
	virtualRouter.HandleFunc("/{cardId}/regenerate", virtualHandler.RegenerateCard).Methods("POST")
	virtualRouter.HandleFunc("/{cardId}/transactions", virtualHandler.GetTransactions).Methods("GET")

	// Card settings routes
	settingsRouter := api.PathPrefix("/cards/settings").Subrouter()
	settingsRouter.HandleFunc("", settingsHandler.GetSettings).Methods("GET")
	settingsRouter.HandleFunc("/default", settingsHandler.UpdateDefaultCards).Methods("PUT")
	settingsRouter.HandleFunc("/security", settingsHandler.UpdateSecuritySettings).Methods("PUT")
	settingsRouter.HandleFunc("/global-limits", settingsHandler.UpdateGlobalLimits).Methods("PUT")
	settingsRouter.HandleFunc("/notifications", settingsHandler.UpdateNotificationSettings).Methods("PUT")
	settingsRouter.HandleFunc("/statement", settingsHandler.UpdateStatementSettings).Methods("PUT")
	settingsRouter.HandleFunc("/pin", settingsHandler.UpdatePINSettings).Methods("PUT")
	settingsRouter.HandleFunc("/authentication", settingsHandler.UpdateAuthenticationSettings).Methods("PUT")

	// Transaction limits routes (works for any card type)
	limitsRouter := api.PathPrefix("/cards/{cardId}/limits").Subrouter()
	limitsRouter.HandleFunc("", limitsHandler.GetLimits).Methods("GET")
	limitsRouter.HandleFunc("/domestic", limitsHandler.UpdateDomesticLimits).Methods("PUT")
	limitsRouter.HandleFunc("/international", limitsHandler.UpdateInternationalLimits).Methods("PUT")

	// Start server
	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", port)
	fmt.Println("API endpoints available at:")
	fmt.Println("  POST   /auth/login")
	fmt.Println("  GET    /api/cards/credit")
	fmt.Println("  GET    /api/cards/debit")
	fmt.Println("  GET    /api/cards/virtual")
	fmt.Println("  ... and many more (see API documentation)")
	
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
