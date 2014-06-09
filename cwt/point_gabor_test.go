package cwt

import (
	"math"
	"math/cmplx"
	"testing"
)

func TestPointGaborCWT(t *testing.T) {
	sampleRate, frameShift := 44100, 441

	// Create Gabor wavelet
	scale := NewCentScale()
	gabor := &PointGaborWavelet{
		SampleRate: sampleRate,
		FrameShift: frameShift,
		Std:        3.0,
		Scale:      scale,
		MinFreq:    scale.MinFreq,
	}

	// Create dummy data
	freq := 200.0
	dummyInput := createSin(freq, sampleRate, 4096)

	// Perform CWT
	complexSpec := gabor.CWT(dummyInput, int(len(dummyInput)/2))

	// Abs
	spec := make([]float64, len(complexSpec))
	for i, val := range complexSpec {
		spec[i] = cmplx.Abs(val)
	}

	maxIndex, _ := maxIndexAndValueOfSlice(spec)
	expectedIndex := scale.FreqToIndex(freq)
	if expectedIndex != maxIndex {
		t.Errorf("The maximum index of CWT coefficients of %f Hz sine data is %v, want %v",
			freq, maxIndex, expectedIndex)
	}

	detectedFreq := scale.IndexToFreq(maxIndex)
	if math.Abs(detectedFreq-freq) > 5.0 {
		t.Errorf("The detected frequency of %f Hz sine data is %f, want about %v",
			freq, detectedFreq, freq)
	}
}

func maxIndexAndValueOfSlice(slice []float64) (int, float64) {
	maxIndex := 0
	maxValue := slice[0]
	for i, val := range slice {
		if val > maxValue {
			maxValue = val
			maxIndex = i
		}
	}
	return maxIndex, maxValue
}

func createSin(freq float64, sampleRate, length int) []float64 {
	sin := make([]float64, length)
	for i := range sin {
		sin[i] = math.Sin(2.0 * math.Pi * freq * float64(i) / float64(sampleRate))
	}
	return sin
}
