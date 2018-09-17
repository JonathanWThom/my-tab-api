package stddrink

import (
	"testing"
	"time"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		percent  float64
		oz       float64
		expected float64
	}{
		{0.12, 5, 1},
		{0.05, 12, 1.00},
		{0.4, 1.5, 1.00},
	}

	for _, test := range tests {
		actual := Calculate(test.percent, test.oz)
		if actual != test.expected {
			t.Errorf("Number of drinks for %f percent and %f oz was incorrect, got: %f, want: %f.",
				test.percent, test.oz, actual, test.expected)
		}
	}
}

func TestTotalStdDrinks(t *testing.T) {
	tests := []struct {
		stddrinks []float64
		expected  float64
	}{
		{[]float64{1}, 1},
		{[]float64{1, 1.5}, 2.5},
	}

	for _, test := range tests {
		actual := TotalStdDrinks(test.stddrinks)
		if actual != test.expected {
			t.Errorf("Total for %v stddrinks was incorrect, got: %f, want: %f.",
				test.stddrinks, actual, test.expected)
		}
	}
}

func TestStddrinksPerDay(t *testing.T) {
	now := time.Now()
	nowMinusOneDay := now.Add(time.Hour * -24)
	nowMinusOneDayOneHour := nowMinusOneDay.Add(time.Hour * -1)

	tests := []struct {
		start     time.Time
		end       time.Time
		stddrinks []float64
		expected  float64
	}{
		{
			nowMinusOneDay,
			now,
			[]float64{1},
			1.0,
		},
		{
			nowMinusOneDay,
			now,
			[]float64{1, 1},
			2.0,
		},
		{
			nowMinusOneDayOneHour,
			now,
			[]float64{1},
			0.96,
		},
	}

	for _, test := range tests {
		actual := StddrinksPerDay(test.start, test.end, test.stddrinks)
		if actual != test.expected {
			t.Errorf("Standard drinks per day for %v start, %v end, and %v stddrinks was incorrect, got: %f, want: %f.",
				test.start, test.end, test.stddrinks, actual, test.expected)
		}
	}
}
