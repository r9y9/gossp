package cwt

import (
	"math"
	"testing"
)

func TestMarrAtTimeZero(t *testing.T) {
	d := &DOG{M: 2}
	expected := 0.8673
	err := math.Abs(d.AtTime(0) - expected)

	tolerance := 1.0e-4
	if err > tolerance {
		t.Errorf("Error %f, want less than %f.", err, tolerance)
	}
}

func TestDOGAtTimeZero(t *testing.T) {
	d := &DOG{M: 6}
	expected := 0.8841
	err := math.Abs(d.AtTime(0) - expected)

	tolerance := 1.0e-4
	if err > tolerance {
		t.Errorf("Error %f, want less than %f.", err, tolerance)
	}
}
