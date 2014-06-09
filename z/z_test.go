package z

import (
	"testing"
)

// A DC removal filter
// https://class.coursera.org/dsp-004/lecture/89
func TestDCCut(t *testing.T) {
	// H(z) = (1- z^-1)/(1 - 0.99z^-1)
	N := 1024
	a := []float64{1, -0.99}
	b := []float64{1, -1}
	amp, _ := SplitFreqZ(FreqZ(b, a, N))

	if amp[0] != 0.0 {
		t.Errorf("%f, the origin must be zero.", amp[0])
	}
	if amp[1] < 0.5 {
		t.Errorf("%f, the amp[1] must be greater than 0.5.", amp[1])
	}

	for i := 5; i < N/2; i++ {
		if amp[i] < 0.9 {
			t.Errorf("%f, want greater than 0.9.", amp[i])
		}
	}
}
