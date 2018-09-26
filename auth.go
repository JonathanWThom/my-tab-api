package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"os"
)

const (
	privKeyPath = "keys/private.key"
	pubKeyPath  = "keys/public.key.pub"
)

var verifyKey *rsa.PublicKey
var signKey *rsa.PrivateKey

func initKeys() {
	var err error
	var signKeyByte []byte
	if os.Getenv("PORT") == "" {
		signKeyByte, err = ioutil.ReadFile(privKeyPath)
	} else {
		signKeyByte = []byte(os.Getenv("PRIVATE_KEY"))
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signKeyByte)
	if err != nil {
		log.Fatalf("[privateKey]: %s\n", err)
	}

	var verifyKeyByte []byte
	if os.Getenv("PORT") == "" {
		verifyKeyByte, err = ioutil.ReadFile(pubKeyPath)
	} else {
		verifyKeyByte = []byte(os.Getenv("PUBLIC_KEY"))
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyKeyByte)
	if err != nil {
		log.Fatalf("[publicKey]: %s\n", err)
		panic(err)
	}
}

type User struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	UUID     string `json:"uuid" db:"uuid"`
}

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Authorization")
}

var userID interface{}

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	setupResponse(&w, r)
	if (*r).Method == http.MethodOptions {
		return
	}

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})

	claims := token.Claims.(jwt.MapClaims)
	userID = claims["userID"].(float64)

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
