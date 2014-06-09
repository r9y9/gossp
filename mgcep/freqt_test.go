package mgcep

import (
	"github.com/r9y9/gossp/3rdparty/sptk"
	"math"
	"testing"
)

func TestFreqTConsistencyWithSPTK(t *testing.T) {
	var (
		sampleRate = 10000
		freq       = 100.0
		bufferSize = 512
		order      = 25
		alpha      = 0.0
	)
	dummyInput := createSin(freq, sampleRate, bufferSize)
	mc := MCep(dummyInput, order, alpha)

	testOrderSet := []int{15, 20, 25}
	testAlphaSet := []float64{0.0, 0.35, 0.41}

	tolerance := 1.0e-64

	for _, a := range testAlphaSet {
		for _, o := range testOrderSet {

			c1 := FreqT(mc, o, a)
			c2 := sptk.FreqT(mc, o, a)

			for i := range c1 {
				err := math.Abs(c1[i] - c2[i])
				if err > tolerance {
					t.Errorf("Error %f, want less than %f.", err, tolerance)
				}
			}
		}
	}
}
