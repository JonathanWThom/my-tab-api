package main

import (
	"encoding/json"
	"net/http"
)

// SignUpHandler is routed from /signup
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	user := User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdUser, err := store.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Do we need to create a json token here?
	JsonResponse(createdUser, w)
}
