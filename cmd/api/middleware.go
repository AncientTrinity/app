package main

import (
	"fmt"
	"net/http"
)

// recoverPanic is middleware that recovers from panics in handlers
// and sends a 500 Internal Server Error response in JSON.
func (a *applicationDependencies) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// defer will be called when the stack unwinds
		defer func() {
			// recover() checks for panics
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				a.serverErrorResponse(w, r, fmt.Errorf("%v", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
