package mgcep

import (
	"math"
)

// C2IR performs conversion from minimum phase cepstral coefficients to
// impule response.
// Note: a bit differ in original function c2ir.
// m+1 = nc
func C2IR(ceps []float64, length int) []float64 {
	h := make([]float64, length)

	m := len(ceps) - 1

	h[0] = math.Exp(ceps[0])
	for n := 1; n < length; n++ {
		d := 0.0
		upperLimit := n
		if n >= m+1 {
			upperLimit = m
		}
		for k := 1; k <= upperLimit; k++ {
			d += float64(k) * ceps[k] * h[n-k]
		}
		h[n] = d / float64(n)
	}

	return h
}
