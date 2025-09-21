package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *applicationDependencies) routes() http.Handler {
	router := httprouter.New()

	// Custom error handlers
	router.NotFound = http.HandlerFunc(a.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

	// Healthcheck
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthcheckHandler)

	// Comments endpoints (CRUD)
	router.HandlerFunc(http.MethodPost, "/v1/comments", a.createCommentHandler)
	router.HandlerFunc(http.MethodGet, "/v1/comments/:id", a.displayCommentHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/comments/:id", a.updateCommentHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/comments/:id", a.deleteCommentHandler)
	router.HandlerFunc(http.MethodGet, "/v1/comments", a.listCommentsHandler) // pagination + sorting

	// Users (Lab 4 – Insert + Get by ID)
	// later you can add POST /v1/users and GET /v1/users/:id handlers
	//router.HandlerFunc(http.MethodPost, "/v1/users", a.createUserHandler)
	//router.HandlerFunc(http.MethodGet, "/v1/users/:id", a.displayUserHandler)

	// Apply middleware in the correct order
	handler := a.recoverPanic(router)
	handler = a.enableCORS(handler)
	handler = a.rateLimit(handler)

	return handler
}
