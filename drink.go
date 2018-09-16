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
