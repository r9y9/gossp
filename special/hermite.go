package special

import (
	"math"
)

// HermitePolynomialsProbAtZero returns Hermite polynomials eveluated at x.
// Reference: http://en.wikipedia.org/wiki/Hermite_polynomials
func HermitePolynomialsProb(x float64, n int) float64 {
	if x == 0.0 {
		return HermitePolynomialsProbAtZero(n)
	}

	y := 0.0

	bound := math.Floor(float64(n) / 2.0)
	for m := 0; m <= int(bound); m++ {
		y += math.Pow(-1, float64(m)) * math.Pow(x, float64(n-2*m)) /
			(float64(factorial(m)) * float64(factorial(n-2*m)) *
				math.Pow(2, float64(m)))
	}

	return y * float64(factorial(n))
}

// HermitePolynomialsProbAtZero returns Hermite polynomials eveluated at zero.
func HermitePolynomialsProbAtZero(n int) float64 {
	if n%2 == 0 {
		return math.Pow(-1, float64(n/2)) * float64(doubleFactorial(n-1))
	} else {
		return 0.0
	}
}

func factorial(x int) int {
	if x == 0 {
		return 1
	}
	return x * factorial(x-1)
}

func doubleFactorial(n int) int {
	prod := 1
	for i := 1; i <= (n+1)/2; i++ {
		prod *= (2*i - 1)
	}
	return prod
}
