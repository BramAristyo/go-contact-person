package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/BramAristyo/rest-api-contact-person/internal/config"
	"github.com/BramAristyo/rest-api-contact-person/internal/database"
	"github.com/BramAristyo/rest-api-contact-person/internal/handler"
	"github.com/BramAristyo/rest-api-contact-person/internal/middleware"
	"github.com/go-playground/validator/v10"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg.DatabaseUrl)
	defer db.Close()

	mux := http.NewServeMux()
	apiMux := http.NewServeMux()

	apiMux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "API is healthy",
		})

		if err != nil {
			return
		}

	})

	validate := validator.New()
	contactHandler := handler.NewContactHandler(db, validate)

	apiMux.HandleFunc("GET /contacts", contactHandler.Paginate)
	apiMux.HandleFunc("GET /contacts/all", contactHandler.GetAll)
	apiMux.HandleFunc("GET /contacts/{id}", contactHandler.GetById)
	apiMux.HandleFunc("POST /contacts", contactHandler.Store)
	apiMux.HandleFunc("PUT /contacts/{id}", contactHandler.Update)
	apiMux.HandleFunc("DELETE /contacts/{id}", contactHandler.Delete)

	mux.Handle("/api/", http.StripPrefix("/api", apiMux))

	// TODO (next steps):
	// 1. Complete Create, Update, and Delete features for contacts.
	// 2. Refactor: Separate logic into service and repository layers.
	// 3. Implement Logging middleware (step 1).
	//
	// Reminder: Continue from here for next development tasks.

	log.Printf("server running on http://localhost:%v", cfg.AppPort)
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, middleware.Logger(middleware.Recovery(mux))))
}
