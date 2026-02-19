package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/BramAristyo/rest-api-contact-person/internal/config"
	"github.com/BramAristyo/rest-api-contact-person/internal/database"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg.DatabaseUrl)
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "API is healthy",
		})

		if err != nil {
			return
		}

	})
	log.Printf("server running on http://localhost:%v", cfg.AppPort)
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, mux))
}
