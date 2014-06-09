package cwt

import (
	"math"
	"testing"
)

func TestPaulAtTimeZero(t *testing.T) {
	p := &Paul{M: 4}
	expected := 1.0789
	err := math.Abs(real(p.AtTime(0)) - expected)

	tolerance := 1.0e-4
	if err > tolerance {
		t.Errorf("Error %f, want less than %f.", err, tolerance)
	}
}
