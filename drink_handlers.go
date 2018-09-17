package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func createDrinkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var drink Drink

	err := json.NewDecoder(r.Body).Decode(&drink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	drink.Percent = drink.Percent / 100

	returnedDrink, err := store.CreateDrink(&drink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonDrink, err := json.Marshal(returnedDrink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonDrink)
}

// getDrinksHandler - GET /drinks
// can receive start and end time params
func getDrinksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	drinks, err := store.GetDrinks(start, end)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	drinkListBytes, err := json.Marshal(drinks)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(drinkListBytes)
}
