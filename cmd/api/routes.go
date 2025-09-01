package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// default error handlers
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// actual routes
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/comments", app.createCommentHandler)
	router.HandlerFunc(http.MethodGet, "/v1/comments/:id", app.displayCommentHandler)


	// wrap with panic recovery
	return app.recoverPanic(router)
	
}
