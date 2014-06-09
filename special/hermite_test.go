package special

import (
	"testing"
)

func TestHermitePolynomialsProbAtZero(t *testing.T) {
	for i := 1; i < 15; i += 2 {
		if v := HermitePolynomialsProbAtZero(i); v != 0.0 {
			t.Errorf("%f, want %f", v, 0.0)
		}
	}

	if v := HermitePolynomialsProbAtZero(0); v != 1.0 {
		t.Errorf("%f, want %f", v, 1.0)
	}
	if v := HermitePolynomialsProbAtZero(2); v != -1.0 {
		t.Errorf("%f, want %f", v, -1.0)
	}
	if v := HermitePolynomialsProbAtZero(4); v != 3.0 {
		t.Errorf("%f, want %f", v, 3.0)
	}
	if v := HermitePolynomialsProbAtZero(6); v != -15.0 {
		t.Errorf("%f, want %f", v, -15.0)
	}
	if v := HermitePolynomialsProbAtZero(8); v != 105.0 {
		t.Errorf("%f, want %f", v, 105.0)
	}
}

func TestHermitePolynomialsProbAtZeroAndHermiteConsistency(t *testing.T) {
	for i := 0; i < 15; i++ {
		v1 := HermitePolynomialsProb(0.0, i)
		v2 := HermitePolynomialsProbAtZero(i)
		if v1 != v2 {
			t.Errorf("%f != %f, want equal.", v1, v2)
		}
	}
}
