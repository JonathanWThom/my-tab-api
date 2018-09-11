package stddrink

import (
	"testing"
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
