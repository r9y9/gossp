package window

import (
	"testing"
)

func testSideValue(t *testing.T, windowedSignal []float64, tolerance float64) {
	if windowedSignal[0] > tolerance {
		t.Errorf("The left side value of windowed signal is %f, want less than %f.", windowedSignal[0], tolerance)
	}
	if windowedSignal[len(windowedSignal)-1] > tolerance {
		t.Errorf("The right side value of windowed signal is %f, want less than %f.", windowedSignal[len(windowedSignal)-1], tolerance)
	}
}

func createUniformVector(length int) []float64 {
	uniform := make([]float64, length)
	for i := range uniform {
		uniform[i] = 1.0
	}
	return uniform
}

func TestBlackman(t *testing.T) {
	length := 512
	dummyInput := createUniformVector(length)

	result := Blackman(dummyInput)
	testSideValue(t, result, 1.0e-3)

	result = BlackmanNormalized(dummyInput)
	testSideValue(t, result, 1.0e-3)
}

func TestHanning(t *testing.T) {
	length := 512
	dummyInput := createUniformVector(length)

	result := Hanning(dummyInput)
	testSideValue(t, result, 1.0e-3)

	result = HanningNormalized(dummyInput)
	testSideValue(t, result, 1.0e-3)
}

func TestHamming(t *testing.T) {
	length := 512
	dummyInput := createUniformVector(length)

	result := Hamming(dummyInput)
	testSideValue(t, result, 1.0e-1)

	result = HammingNormalized(dummyInput)
	testSideValue(t, result, 1.0e-1)
}

func TestGaussian(t *testing.T) {
	length := 512
	dummyInput := createUniformVector(length)

	stddev := 0.4
	result := Gaussian(dummyInput, stddev)
	testSideValue(t, result, 1.0e-1)

	result = GaussianNormalized(dummyInput, stddev)
	testSideValue(t, result, 1.0e-1)
}
