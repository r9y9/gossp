package mgcep

import (
	"math"
)

// IGNorm performs inverse Gamma Normalization to Cepstrum.
func IGNorm(normalizedCeps []float64, gamma float64) []float64 {
	ceps := make([]float64, len(normalizedCeps))

	if gamma == 0.0 {
		copy(ceps, normalizedCeps)
		ceps[0] = math.Log(normalizedCeps[0])
		return ceps
	}

	gain := math.Pow(normalizedCeps[0], gamma)
	for i := range ceps {
		ceps[i] = gain * normalizedCeps[i]
	}
	ceps[0] = (gain - 1.0) / gamma

	return ceps
}
