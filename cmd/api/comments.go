package main

import (
	"errors"
	"fmt"
	"net/http"
	"github.com/yourusername/qod/internal/data"
)

// POST /v1/comments
func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Content string `json:"content"`
		Author  string `json:"author"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := &data.Comment{
		Content: input.Content,
		Author:  input.Author,
	}

	// insert into DB
	err = app.commentModel.Insert(comment)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/comments/%d", comment.ID))

	data := envelope{"comment": comment}
	err = app.writeJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// GET /v1/comments/:id
func (app *application) displayCommentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	comment, err := app.commentModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	data := envelope{"comment": comment}
	err = app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
