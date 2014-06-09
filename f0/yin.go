// The code is a Go port of https://github.com/JorenSix/Pidato (C++).
// Reference:
// H. Kawahara, "YIN, a fundamental frequency estimator for speech and music".
// http://audition.ens.fr/adc/pdf/2002_JASA_YIN.pdf

// Package f0 provides support for fundamental frequency estimation.
package f0

const (
	DefaultBufferSize = 2048
	DefaultThreshold  = 0.15
)

// YIN represents the YIN fundamental frequency estimator.
type YIN struct {
	Buffer     []float64 // Buffer used in the YIN analysis
	BufferSize int
	SampleRate int
	Threshold  float64 // Threshold used in the absolute thresholding step
}

// NewYIN returns a new YIN instantce.
func NewYIN(sampleRate int) *YIN {
	y := new(YIN)
	y.BufferSize = DefaultBufferSize
	y.SampleRate = sampleRate
	y.Buffer = make([]float64, y.BufferSize)
	y.Threshold = DefaultThreshold

	return y
}

// ComputeF0 computes f0 and its confidence for a given audio buffer.
// If no f0 is found, it returns pair of zeros.
func (y *YIN) ComputeF0(audioBuffer []float64) (float64, float64) {
	// set buffer size used in YIN analysis
	if len(audioBuffer) != y.BufferSize {
		y.BufferSize = len(audioBuffer)
		y.Buffer = make([]float64, y.BufferSize)
	}

	// Step. 2
	y.Difference(audioBuffer)

	// Step. 3
	y.CumulativeMeanNormalizedDifference()

	// Step. 4
	tau, probability := y.AbsoluteThreshold()

	// no pitch found
	if probability == 0.0 {
		return 0.0, 0.0
	}

	// Step. 5
	pitchInHz := float64(y.SampleRate) / y.ParabolicInterpolation(tau)

	return pitchInHz, probability
}

// Step 2: Difference function
func (y *YIN) Difference(buffer []float64) {
	delta := 0.0
	for tau := 0; tau < y.BufferSize/2; tau++ {
		for t := 0; t < y.BufferSize/2; t++ {
			delta = buffer[t] - buffer[t+tau]
			y.Buffer[tau] += delta * delta
		}
	}
}

// Step 3: Cumulative mean normalized difference function
func (y *YIN) CumulativeMeanNormalizedDifference() {
	y.Buffer[0] = 1.0
	runningSum := 0.0
	for tau := 1; tau < y.BufferSize/2; tau++ {
		runningSum += y.Buffer[tau]
		y.Buffer[tau] *= float64(tau) / runningSum
	}
}

// Step 4: Absolute threshold
func (y *YIN) AbsoluteThreshold() (int, float64) {
	// first two positions in yinBuffer are always 1
	// So start at the third (index 2)
	for tau := 2; tau < y.BufferSize/2; tau++ {
		if y.Buffer[tau] >= y.Threshold {
			continue
		}
		for {
			if tau+1 >= y.BufferSize/2 || y.Buffer[tau+1] >= y.Buffer[tau] {
				break
			}
			tau++
		}
		// found tau, exit loop and return
		// store the probability
		// From the YIN paper: The threshold determines the list of
		// candidates admitted to the set, and can be interpreted as the
		// proportion of aperiodic power tolerated
		// within a ëëperiodicíí signal.
		//
		// Since we want the periodicity and and not aperiodicity:
		// periodicity = 1 - aperiodicity
		probability := 1.0 - y.Buffer[tau]

		return tau, probability
	}

	// no pitch found
	return 0, 0.0
}

// Step 5: Parabolic Interpolation
func (y *YIN) ParabolicInterpolation(tauEstimate int) float64 {
	betterTau := float64(tauEstimate)
	x0, x2 := 0, 0

	if tauEstimate < 1 {
		x0 = tauEstimate
	} else {
		x0 = tauEstimate - 1
	}

	if tauEstimate+1 < y.BufferSize/2 {
		x2 = tauEstimate + 1
	} else {
		x2 = tauEstimate
	}

	switch {
	case x0 == tauEstimate:
		if y.Buffer[tauEstimate] <= y.Buffer[x2] {
			betterTau = float64(tauEstimate)
		} else {
			betterTau = float64(x2)
		}
	case x2 == tauEstimate:
		if y.Buffer[tauEstimate] <= y.Buffer[x0] {
			betterTau = float64(tauEstimate)
		} else {
			betterTau = float64(x2)
		}
	default:
		s0 := y.Buffer[x0]
		s1 := y.Buffer[tauEstimate]
		s2 := y.Buffer[x2]
		betterTau = float64(tauEstimate) + (s2-s0)/(2.0*(2.0*s1-s2-s0))
	}

	return betterTau
}
