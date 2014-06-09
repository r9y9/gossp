package vocoder

import (
	"github.com/r9y9/gossp/3rdparty/sptk"
	"github.com/r9y9/gossp/mgcep"
	"math"
	"testing"
)

func TestB2MCConsistencyWithSPTK(t *testing.T) {
	var (
		sampleRate = 10000
		freq       = 100.0
		bufferSize = 512
		order      = 25
		alpha      = 0.35
	)
	dummyInput := createSin(freq, sampleRate, bufferSize)
	c := mgcep.MCep(dummyInput, order, alpha)
	b := sptk.MC2B(c, alpha)

	tolerance := 1.0e-64

	mc1 := B2MC(b, alpha)
	mc2 := sptk.B2MC(b, alpha)

	for i := range mc1 {
		err := math.Abs(mc1[i] - mc2[i])
		if err > tolerance {
			t.Errorf("Error %f, want less than %f.", err, tolerance)
		}
	}
}
