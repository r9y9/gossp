package mgcep

import (
	"github.com/r9y9/gossp/3rdparty/sptk"
	"math"
	"testing"
)

func TestGC2GCConsistencyWithSPTK(t *testing.T) {
	var (
		sampleRate = 10000
		freq       = 100.0
		bufferSize = 512
		order      = 25
		alpha      = 0.0
		gamma      = -0.5
	)
	dummyInput := createSin(freq, sampleRate, bufferSize)
	mgc := MGCep(dummyInput, order, alpha, gamma)

	testGammaSet := []float64{0.0, -1.0, -0.75, -0.5, -0.25}
	testOrderSet := []int{15, 20, 25}

	tolerance := 1.0e-64

	for _, g := range testGammaSet {
		for _, o := range testOrderSet {

			c1 := GC2GC(mgc, gamma, o, g)
			c2 := sptk.GC2GC(mgc, gamma, o, g)

			for i := range c1 {
				err := math.Abs(c1[i] - c2[i])
				if err > tolerance {
					t.Errorf("Error %f, want less than %f.", err, tolerance)
				}
			}
		}
	}
}
