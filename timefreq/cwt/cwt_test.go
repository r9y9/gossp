package cwt

import (
	"github.com/r9y9/gossp/io"
	"log"
	"math"
	"testing"
)

var (
	testData []float64
)

func loadReal10kData() []float64 {
	w, werr := io.ReadWav("../../test_files/nanami10k.wav")
	if werr != nil {
		log.Fatal(werr)
	}
	data := w.GetMonoData()
	return data
}

func testWaveletReconstructionBase(t *testing.T, wv *Wavelet,
	data []float64, tolerance float64) {
	spectrogram := wv.CWT(data)
	reconstructed := wv.ICWT(spectrogram)

	meanErr := 0.0
	for i := range reconstructed {
		meanErr += math.Abs(reconstructed[i] - data[i])
	}
	meanErr /= float64(len(reconstructed))
	if meanErr > tolerance {
		t.Errorf("Mean Error %f, want less than %f. ", meanErr, tolerance)
	}
}

func TestMorletReconstructionFactor(t *testing.T) {
	var (
		sampleRate = 10000
		length     = 512
	)

	c := &CentScale{
		MinFreq:      27.5,
		MaxFreq:      5000.0,
		CentInterval: 25,
	}

	dt := 1.0 / float64(sampleRate)
	w := &Wavelet{
		DeltaT: dt,
		Scale:  c,
		Basis:  &Morlet{W0: 6.0},
	}

	C := w.ComputeReconstructionFactor(length)
	expected := 0.776

	tolerance := 0.01
	err := math.Abs(C - expected)
	if err > tolerance {
		t.Errorf("Actual %f, want %f that abs error must be less than %f.",
			C, expected, tolerance)
	}
}

func TestPaulReconstructionFactor(t *testing.T) {
	var (
		sampleRate = 10000
		length     = 512
	)

	c := &CentScale{
		MinFreq:      27.5,
		MaxFreq:      5500.0,
		CentInterval: 25,
	}

	dt := 1.0 / float64(sampleRate)
	w := &Wavelet{
		DeltaT: dt,
		Scale:  c,
		Basis:  &Paul{M: 4},
	}

	C := w.ComputeReconstructionFactor(length)
	expected := 1.132

	tolerance := 0.03
	err := math.Abs(C - expected)
	if err > tolerance {
		t.Errorf("Actual %f, want %f that abs error must be less than %f.",
			C, expected, tolerance)
	}
}

func TestMarrReconstructionFactor(t *testing.T) {
	var (
		sampleRate = 10000
		length     = 512
	)

	c := &CentScale{
		MinFreq:      27.5,
		MaxFreq:      5500.0,
		CentInterval: 25,
	}

	dt := 1.0 / float64(sampleRate)
	w := &Wavelet{
		DeltaT: dt,
		Scale:  c,
		Basis:  &DOG{M: 2},
	}

	C := w.ComputeReconstructionFactor(length)
	expected := 3.541

	tolerance := 0.03
	err := math.Abs(C - expected)
	if err > tolerance {
		t.Errorf("Actual %f, want %f that abs error must be less than %f.",
			C, expected, tolerance)
	}
}

func TestDOGReconstructionFactor(t *testing.T) {
	var (
		sampleRate = 10000
		length     = 512
	)
	dt := 1.0 / float64(sampleRate)

	c := &CentScale{
		MinFreq:      27.5,
		MaxFreq:      5000.0,
		CentInterval: 25,
	}

	w := &Wavelet{
		DeltaT: dt,
		Scale:  c,
		Basis:  &DOG{M: 6},
	}

	C := w.ComputeReconstructionFactor(length)
	expected := 1.966

	tolerance := 0.05
	err := math.Abs(C - expected)
	if err > tolerance {
		t.Errorf("Actual %f, want %f that abs error must be less than %f.",
			C, expected, tolerance)
	}
}

func TestMorletReconstruction(t *testing.T) {
	w0Set := []float64{6.0, 7.0, 8.0, 10.0, 12.0, 14.0, 16.0, 20.0}

	var (
		sampleRate = 10000
		freq       = 200.0
		length     = 512
	)
	data := createSin(freq, sampleRate, length)

	c := &CentScale{
		MinFreq:      27.5,
		MaxFreq:      5000.0,
		CentInterval: 25,
	}

	for _, w0 := range w0Set {
		sampleRate := 10000
		dt := 1.0 / float64(sampleRate)

		w := &Wavelet{
			DeltaT: dt,
			Scale:  c,
			Basis:  &Morlet{W0: w0},
		}

		tolerance := 0.15
		testWaveletReconstructionBase(t, w, data, tolerance)
	}
}

func TestPaulReconstruction(t *testing.T) {
	orderSet := []int{4, 6, 8, 10}

	var (
		sampleRate = 10000
		freq       = 200.0
		length     = 512
	)
	data := createSin(freq, sampleRate, length)

	c := &CentScale{
		MinFreq:      27.5,
		MaxFreq:      5500.0,
		CentInterval: 25,
	}

	for _, m := range orderSet {
		sampleRate := 10000
		dt := 1.0 / float64(sampleRate)
		w := &Wavelet{
			DeltaT: dt,
			Scale:  c,
			Basis:  &Paul{M: m},
		}

		tolerance := 0.15
		testWaveletReconstructionBase(t, w, data, tolerance)
	}
}

func TestDOGReconstruction(t *testing.T) {
	orderSet := []int{2, 4, 6, 8, 10, 20, 30}

	var (
		sampleRate = 10000
		freq       = 200.0
		length     = 512
	)
	data := createSin(freq, sampleRate, length)

	c := &CentScale{
		MinFreq:      27.5,
		MaxFreq:      5000.0,
		CentInterval: 25,
	}

	for _, m := range orderSet {
		sampleRate := 10000
		dt := 1.0 / float64(sampleRate)

		dog := &DOG{M: m}

		w := &Wavelet{
			DeltaT: dt,
			Scale:  c,
			Basis:  dog,
		}

		tolerance := 0.15
		testWaveletReconstructionBase(t, w, data, tolerance)
	}
}

func TestMorletReconstructionForRealData(t *testing.T) {
	var (
		sampleRate = 10000
		data       = loadReal10kData()
	)

	m := &Morlet{W0: 6.0}

	w := &Wavelet{
		DeltaT: 1.0 / float64(sampleRate),
		Scale:  NewCentScale(),
		Basis:  m,
	}

	tolerance := 58.0
	testWaveletReconstructionBase(t, w, data, tolerance)
}

func TestPaulReconstructionForRealData(t *testing.T) {
	var (
		sampleRate = 10000
		data       = loadReal10kData()
	)

	c := &CentScale{
		MinFreq:      27.5,
		MaxFreq:      5500.0,
		CentInterval: 25,
	}

	w := &Wavelet{
		DeltaT: 1.0 / float64(sampleRate),
		Scale:  c,
		Basis:  &Paul{M: 4},
	}

	tolerance := 58.0
	testWaveletReconstructionBase(t, w, data, tolerance)
}

func TestDOGReconstructionForRealData(t *testing.T) {
	var (
		sampleRate = 10000
		data       = loadReal10kData()
	)

	c := &CentScale{
		MinFreq:      27.5,
		MaxFreq:      5000.0,
		CentInterval: 25,
	}

	w := &Wavelet{
		DeltaT: 1.0 / float64(sampleRate),
		Scale:  c,
		Basis:  &DOG{M: 6},
	}

	tolerance := 58.0
	testWaveletReconstructionBase(t, w, data, tolerance)
}
