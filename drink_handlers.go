package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jonathanwthom/my-tab-api/stddrink"
	"net/http"
	"time"
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

	metadata := NewDrinksMetadata()

	if len(drinks) > 0 {
		drinksSlice := Drinks(drinks)
		stddrinkList := drinksSlice.StddrinkList()
		var times []time.Time
		if start == "" || end == "" {
			times = drinksSlice.FirstLastTimes()
		} else {
			times, err = stringsToTimes([]string{start, end})
		}

		if err != nil {
			fmt.Fprintf(w, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		perDay := stddrink.StddrinksPerDay(times[1], times[0], stddrinkList)
		total := stddrink.TotalStdDrinks(stddrinkList)

		metadata.Drinks = drinks
		metadata.StddrinksPerDay = perDay
		metadata.TotalStddrinks = total
	}

	drinkMetadataBytes, err := json.Marshal(metadata)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(drinkMetadataBytes)
}

// DELETE /drinks/{id}
func deleteDrinkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := store.DeleteDrink(id)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
