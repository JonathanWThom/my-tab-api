package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"time"
)

// SignUpHandler is routed from POST /signup
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == http.MethodOptions {
		return
	}

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

	createJWTToken(createdUser, w)
}

// LoginHandler is routed from POST /login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == http.MethodOptions {
		return
	}

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Error in request")
		return
	}

	storedUser, err := store.LoginUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Println("Error logging in")
		fmt.Fprintf(w, "Invalid credentials")
		return
	}

	createJWTToken(storedUser, w)
}

func createJWTToken(user *User, w http.ResponseWriter) {
	signer := jwt.New(jwt.GetSigningMethod("RS256"))

	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1))
	claims["iat"] = time.Now().Unix()
	claims["userID"] = user.ID
	signer.Claims = claims

	tokenString, err := signer.SignedString(signKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		log.Printf("Error signing token: %v\n", err)
	}

	response := Token{tokenString}
	JsonResponse(response, w)
}
