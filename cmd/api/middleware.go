package main

import (
	"fmt"
	"net/http"
)

// recoverPanic is middleware that recovers from panics in handlers
// and sends a 500 Internal Server Error response in JSON.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// defer will be called when the stack unwinds
		defer func() {
			// recover() checks if a panic occurred
			if err := recover(); err != nil {
				// ensure connection is closed
				w.Header().Set("Connection", "close")

				// log & return a safe error response
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		// call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
