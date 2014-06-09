package dct

import (
	"math"
	"testing"
)

func testNearEqual(t *testing.T, x, y []float64, tolerance float64) {
	for i := range x {
		err := math.Abs(x[i] - y[i])
		if err > tolerance {
			t.Errorf("Error %f, want less than %f.", err, tolerance)
		}
	}
}

// DCT-II naive implementation
func dctNaive(x []float64) []float64 {
	dctCoef := make([]float64, len(x))
	N := len(x)
	for k := 0; k < N; k++ {
		sum := 0.0
		for i := 0; i < N; i++ {
			sum += x[i] * math.Cos((float64(i)+0.5)*
				math.Pi*float64(k)/float64(N))
		}
		dctCoef[k] = sum
	}

	return dctCoef
}

// Orthogonal DCT-II naive implementation
func dctOrthogonalNaive(x []float64) []float64 {
	dctCoef := make([]float64, len(x))
	N := len(x)
	for k := 0; k < N; k++ {
		sum := 0.0
		for i := 0; i < N; i++ {
			sum += x[i] * math.Cos((float64(i)+0.5)*
				math.Pi*float64(k)/float64(N))
		}
		dctCoef[k] = sum
		if k == 0 {
			dctCoef[k] *= math.Sqrt(1.0 / float64(N))
		} else {
			dctCoef[k] *= math.Sqrt(2 / float64(N))
		}
	}

	return dctCoef
}

// DCT-III (Inverse DCT-II)
func idctNaive(x []float64) []float64 {
	dctCoef := make([]float64, len(x))
	N := len(x)
	c := 0.5 * x[0]
	for k := 0; k < N; k++ {
		sum := 0.0
		for i := 1; i < N; i++ {
			sum += x[i] * math.Cos(math.Pi*float64(i)/
				float64(N)*(float64(k)+0.5))
		}
		dctCoef[k] = sum + c
	}

	return dctCoef
}

func TestFastDCTAndNaiveDCTConsistency(t *testing.T) {
	testData := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	fastDCTResult := DCT(testData)
	naiveDCTResult := dctNaive(testData)

	if len(fastDCTResult) != len(naiveDCTResult) {
		t.Errorf("Dimention mismatch")
	}

	tolerance := 1.0e-12
	testNearEqual(t, fastDCTResult, naiveDCTResult, tolerance)
}

func TestFastIDCTAndNaiveIDCTConsistency(t *testing.T) {
	testData := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	fastIDCTResult := IDCT(testData)
	naiveIDCTResult := idctNaive(testData)

	if len(fastIDCTResult) != len(naiveIDCTResult) {
		t.Errorf("Dimention mismatch")
	}

	tolerance := 1.0e-12
	testNearEqual(t, fastIDCTResult, naiveIDCTResult, tolerance)
}

func TestFastDCTOrthogonalAndNaiveDCTOrthogonalConsistency(t *testing.T) {
	testData := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	fastDCTResult := DCTOrthogonal(testData)
	naiveDCTResult := dctOrthogonalNaive(testData)

	if len(fastDCTResult) != len(naiveDCTResult) {
		t.Errorf("Dimention mismatch")
	}

	tolerance := 1.0e-12
	testNearEqual(t, fastDCTResult, naiveDCTResult, tolerance)
}

func TestOrthogonalTrasform(t *testing.T) {
	testData := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

	reconstructed := IDCTOrthogonal(DCTOrthogonal(testData))

	tolerance := 1.0e-14
	testNearEqual(t, testData, reconstructed, tolerance)
}
