package sptk

import (
	"fmt"
	"math"
	"testing"
)

// Just tests for calling c-functions

var (
	dummyInput []float64
	sampleRate int
	freq       float64
	bufferSize int
)

func init() {
	sampleRate = 10000
	freq = 100.0
	bufferSize = 512
	dummyInput = createSin(freq, sampleRate, bufferSize)
}

func TestACep(t *testing.T) {
	for i, val := range dummyInput {
		result, _ := ACep(val, 25, 0.98, 0.1, 0.9, 4, 0.0)
		if false {
			fmt.Println(i, result)
		}
	}
}

func TestMCep(t *testing.T) {
	result := MCep(dummyInput, bufferSize,
		25, 0.35, 2, 30, 0.001, 0, 0.0, 0.00001, 0)
	if false {
		for i, val := range result {
			fmt.Println(i, val)
		}
	}
}

func TestMGCep(t *testing.T) {
	gamma := 0.0
	result := MGCepWithDefaultParameters(dummyInput, 25, 0.35, gamma)

	for i, val := range result {
		fmt.Println(i, val)
	}
}

func TestMFCC(t *testing.T) {
	MFCC(dummyInput, sampleRate,
		0.97, 1.0, bufferSize, bufferSize, 12, 20, 22, false, true)
}

func createSin(freq float64, sampleRate, length int) []float64 {
	sin := make([]float64, length)
	for i := range sin {
		sin[i] = math.Sin(2.0 * math.Pi * freq * float64(i) / float64(sampleRate))
	}
	return sin
}
