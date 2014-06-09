package z

import (
	"math"
	"math/cmplx"
)

// FreqZ return frequency response given filter coefficients.
// The name of this function is followed by Matlab.
func FreqZ(b, a []float64, N int) []complex128 {
	H := make([]complex128, N)

	for n := 0; n < N; n++ {
		z := cmplx.Exp(complex(0.0, -2.0*math.Pi*float64(n)/float64(N)))
		numerator, denominator := complex(0, 0), complex(0, 0)
		for i := range b {
			numerator += complex(b[len(b)-1-i], 0.0) *
				cmplx.Pow(z, complex(float64(i), 0))
		}
		for i := range a {
			denominator += complex(a[len(a)-1-i], 0.0) *
				cmplx.Pow(z, complex(float64(i), 0))
		}
		H[n] = numerator / denominator
	}

	return H
}
