package f0

import (
	"math"
	"testing"
)

func testYINBase(t *testing.T, freq float64, sampleRate, bufferSize int,
	toleranceInHz float64) {
	dummyInput := createSin(freq, sampleRate, bufferSize)

	y := NewYIN(sampleRate)
	freqEstimate, _ := y.ComputeF0(dummyInput)
	error := math.Abs(freq - freqEstimate)
	if error > toleranceInHz {
		t.Errorf("The estimate is %f of %f Hz signal, and absolute error is %f, the error want lenn than %f",
			freqEstimate, freq, error, toleranceInHz)
	}
}

func TestYIN(t *testing.T) {
	bufferSizeSet := []int{1000, 2000, 3000, 4000, 5000, 6000, 10000}
	sampleRateSet := []int{48000, 44100, 22050, 16000, 8000}
	freqSet := []float64{100.0, 120.0, 200.0, 400.0}

	// test for all combination
	for _, bufferSize := range bufferSizeSet {
		for _, sampleRate := range sampleRateSet {
			for _, freq := range freqSet {
				testYINBase(t, freq, sampleRate, bufferSize, 1.0)
			}
		}
	}
}

func TestYINAccuracy(t *testing.T) {
	sampleRate, bufferSize := 44100, 4410
	freqTestSet := []float64{100.0, 101.0, 102.0, 103.0, 104.0}
	for _, freq := range freqTestSet {
		testYINBase(t, freq, sampleRate, bufferSize, 1.0e-3)
	}
}

func TestYINLowFrequency(t *testing.T) {
	sampleRate, bufferSize := 44100, 4410
	freqTestSet := []float64{20.0, 30.0, 40.0, 50.0, 60.0, 70.0}

	for _, freq := range freqTestSet {
		testYINBase(t, freq, sampleRate, bufferSize, 1.5)
	}
}

// high in human voice (not instruments)
func TestYINHighFrequency(t *testing.T) {
	sampleRate, bufferSize := 44100, 4410
	freqTestSet := []float64{1000.0, 2000.0, 3000.0}

	for _, freq := range freqTestSet {
		testYINBase(t, freq, sampleRate, bufferSize, 5.0)
	}
}

func TestYINVeryHighFrequency(t *testing.T) {
	sampleRate, bufferSize := 44100, 4410
	freqTestSet := []float64{4500.0, 4600.0, 4700.0, 4800.0, 4900.0, 5000.0}

	for _, freq := range freqTestSet {
		testYINBase(t, freq, sampleRate, bufferSize, 30.0)
	}
}

func createSin(freq float64, sampleRate, length int) []float64 {
	sin := make([]float64, length)
	for i := range sin {
		sin[i] = math.Sin(2.0 * math.Pi * freq * float64(i) / float64(sampleRate))
	}
	return sin
}
