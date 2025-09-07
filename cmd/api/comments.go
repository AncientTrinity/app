package main

import (
	"net/http"

	"github.com/victortillett/app/internal/data"
	"github.com/victortillett/app/internal/validator"
)

func (a *applicationDependencies) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var incomingData struct {
		Content string `json:"content"`
		Author  string `json:"author"`
	}

	err := a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	comment := &data.Comment{
		Content: incomingData.Content,
		Author:  incomingData.Author,
	}

	v := validator.New()
	data.ValidateComment(v, comment)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.commentModel.Insert(comment)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/comments/%d", comment.ID))

	dataResponse := envelope{"comment": comment}
	err = a.writeJSON(w, http.StatusCreated, dataResponse, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) displayCommentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	comment, err := a.commentModel.Get(id)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	dataResponse := envelope{"comment": comment}
	err = a.writeJSON(w, http.StatusOK, dataResponse, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
