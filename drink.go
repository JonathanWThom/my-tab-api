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
	Drinks          []*Drink `json:"drinks"`
	StddrinksPerDay float64  `json:"stddrinks_per_day"`
	TotalStddrinks  float64  `json:"total_stdddrinks"`
}

// Drinks is an array of Drink pointers
type Drinks []*Drink

// StddrinkList returns all the stddrinks of a group of Drinks
func (drinks Drinks) StddrinkList() []float64 {
	var stddrinks []float64
	for _, drink := range drinks {
		stddrinks = append(stddrinks, drink.Stddrink)
	}
	return stddrinks
}

// FirstLastTimes returns the first and last time for a group of Drinks
func (drinks Drinks) FirstLastTimes() []time.Time {
	first := drinks[0].ImbibedOn
	last := drinks[len(drinks)-1].ImbibedOn

	return []time.Time{first, last}
}
