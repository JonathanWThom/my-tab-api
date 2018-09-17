package stddrink

import (
	"time"
)

const stdDrink float64 = 0.6

// Calculate determines who many standard drinks a serving was
func Calculate(percent, oz float64) float64 {
	return truncate((percent * oz) / stdDrink)
}

// StddrinksPerDay determines the average standard drinks per day over a time period
func StddrinksPerDay(start, end time.Time, stddrinks []float64) float64 {
	total := TotalStdDrinks(stddrinks)
	diff := end.Sub(start).Hours() / 24

	return total / float64(diff)
}

// TotalStdDrinks adds standard drinks (or rather, floats)
func TotalStdDrinks(stddrinks []float64) float64 {
	var total float64

	for i := range stddrinks {
		total += stddrinks[i]
	}

	return total
}

func truncate(some float64) float64 {
	return float64(int(some*100)) / 100
}
