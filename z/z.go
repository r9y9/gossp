// Package z provides support for Z transform to analyze digital filters.
package z

import (
	"math"
	"math/cmplx"
)

// SplitFreqZ splits a frequency response to amplitude and phase response.
func SplitFreqZ(freqz []complex128) ([]float64, []float64) {
	N := len(freqz)

	ampResponse := make([]float64, N)
	phaseResponse := make([]float64, N)

	for i, val := range freqz {
		ampResponse[i] = cmplx.Abs(val)
		phaseResponse[i] = angle(val)
	}

	return ampResponse, phaseResponse
}

// Angle returns angle of complex number in radian.
func angle(z complex128) float64 {
	return math.Atan2(imag(z), real(z))
}
