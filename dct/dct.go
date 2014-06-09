// Package dct provides support for Discrete Cosine Transform (DCT).
package dct

import (
	"github.com/mjibson/go-dsp/fft"
	"math"
	"math/cmplx"
)

// References:
//     http://www.lokminglui.com/DCT_TR802.pdf
//     http://fourier.eng.hmc.edu/e161/lectures/dct/node2.html

// DCT returns DCT-II coefficients given an input signal.
func DCT(x []float64) []float64 {
	n := len(x)
	nh := n / 2

	evenVector := make([]float64, n)
	for i := 0; i < nh; i++ {
		evenVector[i], evenVector[n-1-i] = x[2*i], x[2*i+1]
	}

	array := fft.FFTReal(evenVector)
	theta := math.Pi / (2.0 * float64(n))

	// 1/4 Shift
	for k := 1; k < nh; k++ {
		w := cmplx.Exp(complex(0, -float64(k)*theta))
		wCont := -complex(imag(w), real(w))
		array[k] *= w
		array[n-k] *= wCont
	}
	array[nh] *= complex(math.Cos(theta*float64(nh)), 0.0)

	dctCoef := make([]float64, n)
	for i := range dctCoef {
		dctCoef[i] = real(array[i])
	}

	return dctCoef
}

// DCTOrthogonal returns orthogonal DCT-II coefficients given an input signal.
func DCTOrthogonal(x []float64) []float64 {
	dctCoef := DCT(x)
	n := len(x)

	c0 := math.Sqrt(1.0 / float64(n))
	c1 := math.Sqrt(2.0 / float64(n))

	for i := 1; i < n; i++ {
		dctCoef[i] *= c1
	}
	dctCoef[0] *= c0

	return dctCoef
}

// IDCT returns Inverse DCT (DCT-III) given a input signal.
func IDCT(dctCoef []float64) []float64 {
	n := len(dctCoef)
	nh := n / 2
	array := make([]complex128, n)

	theta := math.Pi / (2.0 * float64(n))

	// -1/4 Shift
	for k := 1; k < nh; k++ {
		w := cmplx.Exp(complex(0, float64(k)*theta))
		wCont := complex(imag(w), real(w))
		array[k] = w * complex(dctCoef[k], 0)
		array[n-k] = wCont * complex(dctCoef[n-k], 0)
	}
	array[0] = complex(dctCoef[0], 0)
	array[nh] = complex(dctCoef[nh], 0) *
		cmplx.Exp(complex(0, float64(nh)*theta))

	y := fft.IFFT(array)

	x := make([]float64, n)
	for i := 0; i < nh; i++ {
		x[2*i] = float64(n) * real(y[i])
		x[2*i+1] = float64(n) * real(y[n-1-i])
	}

	return x
}

// IDCTOrthogonal returns Inverse DCT (DCT-III) coefficients given a input
// signal.
func IDCTOrthogonal(dctCoef []float64) []float64 {
	n := len(dctCoef)
	nh := n / 2
	array := make([]complex128, n)

	theta := math.Pi / (2.0 * float64(n))

	c0 := math.Sqrt(1.0 / float64(n))
	c1 := math.Sqrt(2.0 / float64(n))

	// -1/4 Shift
	for k := 1; k < nh; k++ {
		w := cmplx.Exp(complex(0, float64(k)*theta))
		wCont := complex(imag(w), real(w))
		array[k] = w * complex(c1*dctCoef[k], 0)
		array[n-k] = wCont * complex(c1*dctCoef[n-k], 0)
	}
	array[0] = complex(c0*dctCoef[0], 0)
	array[nh] = complex(c1*dctCoef[nh], 0) *
		cmplx.Exp(complex(0, float64(nh)*theta))

	y := fft.IFFT(array)

	x := make([]float64, n)
	for i := 0; i < nh; i++ {
		x[2*i] = float64(n) * real(y[i])
		x[2*i+1] = float64(n) * real(y[n-1-i])
	}

	return x
}
