package mgcep

import (
	"math"
)

// GNorm performs Gamma Normalization to Cepstrum.
func GNorm(ceps []float64, gamma float64) []float64 {
	normalizedCeps := make([]float64, len(ceps))

	if gamma == 0.0 {
		copy(normalizedCeps, ceps)
		normalizedCeps[0] = math.Exp(ceps[0])
		return normalizedCeps
	} else {
		gain := 1.0 + gamma*ceps[0]
		m := len(ceps) - 1
		for i := m; i >= 1; i-- {
			normalizedCeps[i] = ceps[i] / gain
		}
		normalizedCeps[0] = math.Pow(gain, 1.0/gamma)
		return normalizedCeps
	}
}
