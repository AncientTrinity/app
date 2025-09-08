package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	//"victortillett.net/basic/internal/data"
	//"github.com/julienschmidt/httprouter"
)

func (a *applicationDependencies) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	// create a struct to hold a comment
	// struct tags (`json:"..."`) ensure JSON keys are lowercase
	var incomingData struct {
		Content string `json:"content"`
		Author  string `json:"author"`
	}

	// decode JSON request body
	err := json.NewDecoder(r.Body).Decode(&incomingData)
	if err != nil {
		a.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// for now, display the result
	fmt.Fprintf(w, "%+v\n", incomingData)
}