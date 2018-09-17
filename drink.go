package main

import "time"

type Drink struct {
	ID        int       `json:"id"`
	Percent   float64   `json:"percent,string"`
	Oz        float64   `json:"oz,string"`
	Stddrink  float64   `json:"stddrink"`
	ImbibedOn time.Time `json:"imbibedOn"`
	UserID    int       `json:"user_id"`
}

// DrinksMetadata stores an array of drinks plus cumulative metadata about all returned drinks
type DrinksMetadata struct {
	Drinks          []Drink `json:"drinks"`
	StddrinksPerDay float64 `json:"stddrinks_per_day"`
	TotalStddrinks  float64 `json:"total_stdddrinks"`
}
