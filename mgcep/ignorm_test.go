package mgcep

import (
	"github.com/r9y9/gossp/3rdparty/sptk"
	"math"
	"testing"
)

func TestIGNormConsistencyWithSPTK(t *testing.T) {
	var (
		sampleRate = 10000
		freq       = 100.0
		bufferSize = 512
		order      = 25
		alpha      = 0.35
	)
	dummyInput := createSin(freq, sampleRate, bufferSize)

	gamma := -0.5
	mgc := MGCep(dummyInput, order, alpha, gamma)

	testGammaSet := []float64{0.0, -1.0, -0.75, -0.5, -0.25}

	tolerance := 1.0e-15

	for _, g := range testGammaSet {
		c1 := IGNorm(mgc, g)
		c2 := sptk.IGNorm(mgc, g)

		for i := range c1 {
			err := math.Abs(c1[i] - c2[i])
			if err > tolerance {
				t.Errorf("Error %f, want less than %f.", err, tolerance)
			}
		}
	}
}
