package cwt

import (
	"math"
	"testing"
)

func TestMorletAtTimeZero(t *testing.T) {
	m := &Morlet{W0: 6.0}
	expected := 0.7511
	err := math.Abs(m.AtTimeZero() - expected)

	tolerance := 1.0e-4
	if err > tolerance {
		t.Errorf("Error %f, want less than %f.", err, tolerance)
	}
}
