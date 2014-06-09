package mgcep

import (
	"github.com/r9y9/gossp/3rdparty/sptk"
	"math"
	"testing"
)

func TestC2IRConsistencyWithSPTK(t *testing.T) {
	var (
		sampleRate = 10000
		freq       = 100.0
		bufferSize = 512
		order      = 25
		alpha      = 0.0
		length     = 512
	)
	dummyInput := createSin(freq, sampleRate, bufferSize)
	c := MCep(dummyInput, order, alpha)

	tolerance := 1.0e-64

	c1 := C2IR(c, length)
	c2 := sptk.C2IR(c, length)

	for i := range c1 {
		err := math.Abs(c1[i] - c2[i])
		if err > tolerance {
			t.Errorf("Error %f, want less than %f.", err, tolerance)
		}
	}
}
