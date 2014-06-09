package vocoder

import (
	"github.com/r9y9/gossp/3rdparty/sptk"
	"github.com/r9y9/gossp/mgcep"
	"math"
	"testing"
)

func TestMC2B(t *testing.T) {
	var (
		sampleRate = 10000
		freq       = 100.0
		bufferSize = 512
		order      = 25
		alpha      = 0.35
	)
	dummyInput := createSin(freq, sampleRate, bufferSize)
	mc := mgcep.MCep(dummyInput, order, alpha)

	b1 := MC2B(mc, alpha)
	b2 := sptk.MC2B(mc, alpha)

	tolerance := 1.0e-64

	for i := range b1 {
		err := math.Abs(b1[i] - b2[i])
		if err > tolerance {
			t.Errorf("Error %f, want less than %f.", err, tolerance)
		}
	}

}
