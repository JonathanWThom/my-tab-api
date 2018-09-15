package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

const (
	privKeyPath = "keys/private.key"
	pubKeyPath  = "keys/public.key.pub"
)

var verifyKey *rsa.PublicKey
var signKey *rsa.PrivateKey

func initKeys() {
	var err error
	signKeyByte, err := ioutil.ReadFile(privKeyPath)
	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signKeyByte)
	if err != nil {
		log.Fatalf("[privateKey]: %s\n", err)
	}

	verifyKeyByte, err := ioutil.ReadFile(pubKeyPath)
	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyKeyByte)
	if err != nil {
		log.Fatalf("[publicyKey]: %s\n", err)
		panic(err)
	}
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{"Gained access to protected resource"}
	JsonResponse(response, w)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user UserCredentials

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Error in request")
		return
	}

	if strings.ToLower(user.Username) != "jonathan" {
		if user.Password != "123123" {
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("Error logging in")
			fmt.Fprintf(w, "Invalid credentials")
			return
		}
	}

	signer := jwt.New(jwt.GetSigningMethod("RS256"))

	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1))
	claims["iat"] = time.Now().Unix()
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

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})

	if err == nil {
		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorised access to this resource")
	}
}

func JsonResponse(response interface{}, w http.ResponseWriter) {
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
