package cwt

import (
	"errors"
	"math"
)

// Morlet represents the Morlet wavelet function.
type Morlet struct {
	W0 float64 // frequency
}

// FrequencyRep returns wavelet in frequency domain.
func (m *Morlet) FrequencyRep(omega []float64, scale float64) []float64 {
	rep := make([]float64, len(omega))

	sw := make([]float64, len(omega))
	for i := range sw {
		sw[i] = omega[i] * scale
	}

	norm := math.Pow(math.Pi, -0.25)
	for i := range rep {
		if omega[i] > 0 {
			rep[i] = norm * math.Exp(-0.5*math.Pow(sw[i]-m.W0, 2.0))
		} else {
			rep[i] = 0.0
		}
	}

	return rep
}

func (m *Morlet) AtTimeZero() float64 {
	return math.Pow(math.Pi, -0.25)
}

func (m *Morlet) ReconstructionFactor() (float64, error) {
	switch m.W0 {
	case 5.0:
		return 0.9484, nil
	case 5.336:
		return 0.8831, nil
	case 6.0:
		return 0.776, nil
	case 7.0:
		return 0.6616, nil
	case 8.0:
		return 0.5758, nil
	case 10.0:
		return 0.4579, nil
	case 12.0:
		return 0.3804, nil
	case 14.0:
		return 0.3254, nil
	case 16.0:
		return 0.2844, nil
	case 20.0:
		return 0.2272, nil
	default:
		return 0.0, errors.New("not defined")
	}
}

func (m *Morlet) FourierPeriod(scale float64) float64 {
	return 4.0 * float64(math.Pi) * scale /
		(m.W0 + math.Sqrt(2.0+m.W0*m.W0))
}
