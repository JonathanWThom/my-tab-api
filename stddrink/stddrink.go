package stddrink

const Drink float64 = 0.6

func Calculate(percent, oz float64) float64 {
	return truncate((percent * oz) / Drink)
}

func truncate(some float64) float64 {
	return float64(int(some*100)) / 100
}
