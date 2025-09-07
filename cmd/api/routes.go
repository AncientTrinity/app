package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *alication) routes() http.Handler {
	router := httprouter.New()

	// default error handlers
	router.NotFound = http.HandlerFunc(a.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

	// actual routes
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/comments", a.createCommentHandler)
	router.HandlerFunc(http.MethodGet, "/v1/comments/:id", a.displayCommentHandler)


	// wrap with panic recovery
	//return a.recoverPanic(router)

	// Return the httprouter instance.
	return router
	
}
