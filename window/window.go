package window

import (
	"math"
)

// Windowing returns windowed signal for given input and window signal.
func Windowing(input, window []float64) []float64 {
	result := make([]float64, len(input))

	for i := range input {
		result[i] = input[i] * window[i]
	}

	return result
}

// WindowingNormalized returns power-normalized windowed signal given
// input and window signal.
func WindowingNormalized(input, window []float64) []float64 {
	result := make([]float64, len(input))

	// Compute power of window
	sum := 0.0
	for _, val := range window {
		sum += val
	}
	power := math.Sqrt(sum)

	for i := range input {
		result[i] = input[i] * (window[i] / power)
	}

	return result
}

func Blackman(input []float64) []float64 {
	return Windowing(input, CreateBlackman(len(input)))
}

func BlackmanNormalized(input []float64) []float64 {
	return WindowingNormalized(input, CreateBlackman(len(input)))
}

func Hamming(input []float64) []float64 {
	return Windowing(input, CreateHamming(len(input)))
}

func HammingNormalized(input []float64) []float64 {
	return WindowingNormalized(input, CreateHamming(len(input)))
}

func Hanning(input []float64) []float64 {
	return Windowing(input, CreateHanning(len(input)))
}

func HanningNormalized(input []float64) []float64 {
	return WindowingNormalized(input, CreateHanning(len(input)))
}

func Gaussian(input []float64, stddev float64) []float64 {
	return Windowing(input, CreateGaussian(len(input), stddev))
}

func GaussianNormalized(input []float64, stddev float64) []float64 {
	return WindowingNormalized(input, CreateGaussian(len(input), stddev))
}
