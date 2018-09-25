package main

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Starting server...")
	addr, err := determineListenAddress()
	if err != nil {
		panic(err)
	}
	initKeys()
	connServer := determineDbUrl()
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
	)).Methods("GET", "OPTIONS")

	router.Handle("/drinks", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(createDrinkHandler)),
	)).Methods("POST", "OPTIONS")

	router.HandleFunc("/signup", SignUpHandler).Methods("POST")
	router.HandleFunc("/login", LoginHandler).Methods("POST")

	fmt.Println("Now serving on port 8000")

	http.ListenAndServe(addr, router)
}

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	return ":" + port, nil
}

func determineDbUrl() string {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return "dbname=my_tab sslmode=disable"
	} else {
		return url
	}
}
