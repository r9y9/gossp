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

func TestMGLSAConsistencyWithSPTK(t *testing.T) {
	var (
		testData []float64
		order    = 25
		alpha    = 0.41
		stage    = 12
		gamma    = -1.0 / float64(stage)
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

	mgc := mgcep.MGCep(window.BlackmanNormalized(data), order, alpha, gamma)
	filterCoef := MGCep2MGLSAFilterCoef(mgc, alpha, gamma)

	mf := NewMGLSAFilter(order, alpha, stage)

	// tricky allocation based on the SPTK (just for test)
	d := make([]float64, (order+1)*stage)

	tolerance := 1.0e-64

	for _, val := range data {
		y1 := mf.Filter(val, filterCoef)
		y2 := sptk.MGLSADF(val, filterCoef, order, alpha, stage, d)
		err := math.Abs(y1 - y2)
		if err > tolerance {
			t.Errorf("Error: %f, want less than %f.", err, tolerance)
		}
	}

}
