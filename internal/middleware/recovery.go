package middleware

import (
	"log"
	"net/http"

	"github.com/BramAristyo/rest-api-contact-person/pkg/response"
)

/*
Return http.Handler because we want to chain it with other middlewares,
and it takes next http.Handler as parameter to call the next handler in
the chain after recovery is done.

Is not linear like express.js where you call next() to move to the next middleware,
but we wrap the next handler inside our middleware function, so we have control
over when to call the next handler, and we can do things before and after it.
*/
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
