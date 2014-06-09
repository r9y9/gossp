package mgcep

import (
	"math"
	"testing"
)

// MGCep analysis if gamma = 0, MGCep analysis is corresponds to MCep analysis
func TestMCepsAsSpecialCaseOfMGCep(t *testing.T) {
	var (
		sampleRate = 10000
		freq       = 100.0
		bufferSize = 512
		order      = 20
		alpha      = 0.35
		gamma      = 0.0
	)
	dummyInput := createSin(freq, sampleRate, bufferSize)

	c1 := MCep(dummyInput, order, alpha)
	c2 := MGCep(dummyInput, order, alpha, gamma)

	tolerance := 1.0e-6

	for i := range c1 {
		err := math.Abs(c1[i] - c2[i])
		if err > tolerance {
			t.Errorf("Error %f at index %d, want less than %f.",
				err, i, tolerance)
		}
	}
}

// MCep analysis if alpha = 0, MCep analysis is corresponds to
// Unbiased Estimation of log spectrum.
func TestUELSAsSpecialCaseOfMGCep(t *testing.T) {
	var (
		sampleRate = 10000
		freq       = 100.0
		bufferSize = 512
		order      = 25
		alpha      = 0.0
		gamma      = 0.0
	)
	dummyInput := createSin(freq, sampleRate, bufferSize)

	c1 := UELS(dummyInput, order)
	c2 := MGCep(dummyInput, order, alpha, gamma)

	tolerance := 1.0e-5

	// TODO(ryuichi): check 0th order consistency
	for i := 1; i < len(c1); i++ {
		err := math.Abs(c1[i] - c2[i])
		if err > tolerance {
			t.Errorf("Error %f at index %d, want less than %f.",
				err, i, tolerance)
		}
	}
}

func createSin(freq float64, sampleRate, length int) []float64 {
	sin := make([]float64, length)
	for i := range sin {
		sin[i] = math.Sin(2.0 * math.Pi * freq * float64(i) / float64(sampleRate))
	}
	return sin
}
