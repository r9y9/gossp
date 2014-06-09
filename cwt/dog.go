package cwt

import (
	"errors"
	"github.com/r9y9/gossp/special"
	"math"
)

// DOG represents the Derivative of Gaussian (DOG) wavelet function.
type DOG struct {
	M int // M-th order derivative
}

// FrequencyRep returns wavelet in frequency domain.
func (d *DOG) FrequencyRep(omega []float64, scale float64) []float64 {
	rep := make([]float64, len(omega))

	sw := make([]float64, len(omega))
	for i := range sw {
		sw[i] = omega[i] * scale
	}

	norm := -math.Pow(-1, math.Floor(float64(d.M/2.0))) /
		math.Sqrt(math.Gamma(float64(d.M)+0.5))
	for i := range rep {
		rep[i] = norm * math.Pow(sw[i], float64(d.M)) *
			math.Exp(-0.5*sw[i]*sw[i])
	}

	return rep
}

// TODO(ryuichi) write test
func (d *DOG) AtTime(t float64) float64 {
	norm := math.Pow(-1, float64(d.M+1)) /
		math.Sqrt(math.Gamma(float64(d.M)+0.5))

	return norm * special.HermitePolynomialsProb(t, d.M) * math.Exp(-0.5*t*t)
}

func (d *DOG) AtTimeZero() float64 {
	return d.AtTime(0)
}

func (d *DOG) ReconstructionFactor() (float64, error) {
	switch d.M {
	case 2:
		return 3.5987, nil
	case 4:
		return 2.4014, nil
	case 6:
		return 1.9212, nil
	case 8:
		return 1.6467, nil
	case 12:
		return 1.3307, nil
	case 16:
		return 1.1464, nil
	case 20:
		return 1.0222, nil
	case 30:
		return 0.8312, nil
	case 40:
		return 0.7183, nil
	case 60:
		return 0.5853, nil
	default:
		return 0.0, errors.New("not defined")
	}
}

func (d *DOG) FourierPeriod(scale float64) float64 {
	return 2.0 * float64(math.Pi) * scale / math.Sqrt(float64(d.M)+0.5)
}
