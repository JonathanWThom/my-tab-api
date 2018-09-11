package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
)

func main() {
	fmt.Println("Starting server...")
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
	router.HandleFunc("/drinks", getDrinksHandler).Methods("GET")
	router.HandleFunc("/drinks", createDrinkHandler).Methods("POST")
	fmt.Println("Now serving on port 8000")

	http.ListenAndServe(":8000", router)
}
