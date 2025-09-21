package main

import (
	"fmt"
	"net/http"

	"victortillett.net/basic/internal/data"
	"victortillett.net/basic/internal/validator"
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
		switch {
		case err == data.ErrRecordNotFound:
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	dataResponse := envelope{"comment": comment}
	err = a.writeJSON(w, http.StatusOK, dataResponse, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}


func (a *applicationDependencies) updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	comment, err := a.commentModel.Get(id)
	if err != nil {
		switch {
		case err == data.ErrRecordNotFound:
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// Input can be partial
	var incomingData struct {
		Content *string `json:"content"`
		Author  *string `json:"author"`
	}

	err = a.readJSON(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	if incomingData.Content != nil {
		comment.Content = *incomingData.Content
	}
	if incomingData.Author != nil {
		comment.Author = *incomingData.Author
	}

	v := validator.New()
	data.ValidateComment(v, comment)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = a.commentModel.Update(comment)
	if err != nil {
		switch {
		case err == data.ErrEditConflict:
			a.editConflictResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	dataResponse := envelope{"comment": comment}
	err = a.writeJSON(w, http.StatusOK, dataResponse, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

func (a *applicationDependencies) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.commentModel.Delete(id)
	if err != nil {
		switch {
		case err == data.ErrRecordNotFound:
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	err = a.writeJSON(w, http.StatusOK, envelope{"message": "comment successfully deleted"}, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// for pagination
func (a *applicationDependencies) listCommentsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page := a.readInt(query, "page", 1)
	pageSize := a.readInt(query, "page_size", 10)
	sort := query.Get("sort")
	if sort == "" {
		sort = "id"
	}

	comments, metadata, err := a.commentModel.GetAll(page, pageSize, sort)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	dataResponse := envelope{
		"comments": comments,
		"metadata": metadata,
	}
	err = a.writeJSON(w, http.StatusOK, dataResponse, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
