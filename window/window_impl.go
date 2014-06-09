package window

import (
	"math"
)

// CreateBlackman returns blackman window function.
func CreateBlackman(length int) []float64 {
	blackman := make([]float64, length)

	arg := 2.0 * math.Pi / float64(length-1)
	for i := range blackman {
		blackman[i] = 0.42 - 0.5*math.Cos(arg*float64(i)) + 0.08*math.Cos(2.0*arg*float64(i))
	}

	return blackman
}

// CreateHanning returns hanning window function.
func CreateHanning(length int) []float64 {
	hanning := make([]float64, length)

	arg := 2.0 * math.Pi / float64(length-1)
	for i := range hanning {
		hanning[i] = 0.5 - 0.5*math.Cos(arg*float64(i))
	}

	return hanning
}

// CreateHamming returns hamming window function.
func CreateHamming(length int) []float64 {
	hamming := make([]float64, length)

	arg := 2.0 * math.Pi / float64(length-1)
	for i := range hamming {
		hamming[i] = 0.54 - 0.46*math.Cos(arg*float64(i))
	}

	return hamming
}

// CreateGaussian returns Gaussian window function.
// It is recommended that stddev (standard deviation) is set to around 0.4.
func CreateGaussian(length int, stddev float64) []float64 {
	gaussian := make([]float64, length)

	for i := range gaussian {
		coef := (float64(i) - float64(length-1)/2.0) / (stddev * float64(length-1) / 2.0)
		gaussian[i] = math.Exp(-0.5 * coef * coef)
	}

	return gaussian
}
