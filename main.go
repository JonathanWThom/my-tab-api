package main

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
)

func main() {
	fmt.Println("Starting server...")
	initKeys()
	connServer := "dbname=my_tab sslmode=disable"
	db, err := sql.Open("postgres", connServer)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	InitStore(&dbStore{db: db})
	router := mux.NewRouter()

	router.Handle("/drinks", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(getDrinksHandler)),
	)).Methods("GET")

	router.Handle("/drinks", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(getDrinksHandler)),
	)).Methods("POST")

	router.HandleFunc("/login", LoginHandler).Methods("POST")

	fmt.Println("Now serving on port 8000")

	http.ListenAndServe(":8000", router)
}
