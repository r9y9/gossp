package vocoder

import (
	"github.com/r9y9/go-dsp/wav"
	"github.com/r9y9/gossp/3rdparty/sptk"
	"github.com/r9y9/gossp/mgcep"
	"github.com/r9y9/gossp/window"
	"log"
	"math"
	"os"
	"testing"
)

func TestMLSAConsistencyWithSPTK(t *testing.T) {
	var (
		testData []float64
		order    = 24
		alpha    = 0.41
		pd       = 5
	)

	file, err := os.Open("../test_files/test16k.wav")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	w, werr := wav.ReadWav(file)
	if werr != nil {
		log.Fatal(werr)
	}
	testData = w.GetMonoData()

	data := testData[:512]

	mc := mgcep.MCep(window.BlackmanNormalized(data), order, alpha)
	filterCoef := MCep2MLSAFilterCoef(mc, alpha)

	mf := NewMLSAFilter(order, alpha, pd)

	// tricky allocation based on the SPTK (just for test)
	d := make([]float64, 3*(pd+1)+pd*(order+2))

	tolerance := 1.0e-10

	for _, val := range data {
		y1 := mf.Filter(val, filterCoef)
		y2 := sptk.MLSADF(val, filterCoef, order, alpha, pd, d)
		err := math.Abs(y1 - y2)
		if err > tolerance {
			t.Errorf("Error: %f, want less than %f.", err, tolerance)
		}
	}
}
