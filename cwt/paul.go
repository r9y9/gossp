package cwt

import (
	"errors"
	"math"
)

// Paul represents the Paul wavelet function.
type Paul struct {
	M int // order
}

// FrequencyRep returns wavelet in frequency domain.
func (p *Paul) FrequencyRep(omega []float64, scale float64) []float64 {
	rep := make([]float64, len(omega))

	sw := make([]float64, len(omega))
	for i := range sw {
		sw[i] = omega[i] * scale
	}

	norm := math.Pow(2.0, float64(p.M)) /
		math.Sqrt(float64(p.M)*float64(factorial((2*p.M-1))))
	for i := range rep {
		if omega[i] > 0 {
			rep[i] = norm * math.Pow(sw[i], float64(p.M)) *
				math.Exp(-sw[i])
		} else {
			rep[i] = 0.0
		}
	}

	return rep
}

func (p *Paul) AtTime(t float64) complex128 {
	// Real part
	r := math.Pow(2.0, float64(p.M)) * math.Pow(-1, math.Floor(float64(p.M/2))) *
		float64(factorial(p.M)) /
		math.Sqrt(float64(math.Pi)*float64(factorial(2*p.M)))

	if p.M%2 == 0 {
		return complex(r, 0)
	} else {
		return complex(0, r)
	}
}

func (p *Paul) AtTimeZero() float64 {
	if p.M%2 == 0 {
		return real(p.AtTime(0))
	} else {
		// TODO(ryuichi)
		panic("Paul: Order must be power of 2.")
	}
}

func (p *Paul) ReconstructionFactor() (float64, error) {
	switch p.M {
	case 4:
		return 1.133, nil
	case 6:
		return -0.7554, nil
	case 8:
		return 0.55665, nil
	case 10:
		return -0.4532, nil
	case 16:
		return 0.2833, nil
	case 20:
		return 0.2266, nil
	case 30:
		return -0.1511, nil
	case 40:
		return 0.1133, nil
	default:
		return 0.0, errors.New("not defined")
	}
}

func (p *Paul) FourierPeriod(scale float64) float64 {
	return 4.0 * float64(math.Pi) * scale / float64(2*p.M+1.0)
}

func factorial(x int) int {
	if x == 0 {
		return 1
	}
	return x * factorial(x-1)
}
