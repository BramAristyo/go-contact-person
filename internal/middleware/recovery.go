package middleware

import (
	"log"
	"net/http"

	"github.com/BramAristyo/rest-api-contact-person/pkg/response"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				response.WriteError(w, "Internal Server error", 500)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}
