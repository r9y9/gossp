package excite

import (
	"math"
	"math/rand"
)

// PulseExcite represents generating excitation signals.
// if UseGauss is set to true, gaussian random samples are generated for
// non-voice segments, pseudo random samples otherwise.
type PulseExcite struct {
	SampleRate int
	FrameShift int
	UseGauss   bool

	// Number of samples after pulse generated is used to keep continuity
	// of successive two frames.
	numSamplesAfterPulseGenerated int
}

// NewPulseExcite returns its instance with sample rate and framshift.
func NewPulseExcite(sampleRate, frameShift int) *PulseExcite {
	return &PulseExcite{sampleRate, frameShift, false, 0}
}

// Gerenate generates an excitation signal from f0 sequence. If the
// unvoiced segment is detected (segment of zero f0), this generates
// gaussian or pseudo random samples, f0-dependent excitation otherwise.
func (e *PulseExcite) Generate(f0Sequence []float64) []float64 {
	excite := make([]float64, e.FrameShift*len(f0Sequence))

	// reset
	e.numSamplesAfterPulseGenerated = 0
	previousF0 := f0Sequence[0]
	for i, currentF0 := range f0Sequence {
		if i > 0 {
			previousF0 = f0Sequence[i-1]
		}
		exciteForFrame := e.GenerateOneFrame(currentF0, previousF0)
		for j := 0; j < e.FrameShift; j++ {
			excite[i*e.FrameShift+j] = exciteForFrame[j]
		}
	}

	return excite
}

// GenerateOneFrame generates an excitation signal from successive two f0.
// If the given f0 have zero value(s), GenerateOneFrame generates random
// samples, f0-dependent excitation otherwise.
func (e *PulseExcite) GenerateOneFrame(f01, f02 float64) []float64 {
	excite := make([]float64, e.FrameShift)

	// Generate random samples
	if f01 == 0.0 || f02 == 0.0 {
		e.numSamplesAfterPulseGenerated = 0
		for i := range excite {
			if e.UseGauss {
				excite[i] = rand.NormFloat64()
			} else {
				excite[i] = rand.Float64()
			}
		}
		return excite
	}

	// Normalize by samplerate
	normalizedF01 := float64(e.SampleRate) / f01
	normalizedF02 := float64(e.SampleRate) / f02
	slope := (normalizedF02 - normalizedF01) / float64(e.FrameShift)

	// Generate f0-dependent excitation
	for i := 0; i < e.FrameShift; i++ {
		// f0 is linearly interpolated
		linearlyInterpolatedF0 := normalizedF01 + slope*float64(i)

		if e.numSamplesAfterPulseGenerated > int(linearlyInterpolatedF0) {
			// Generate Pulse
			excite[i] = math.Sqrt(linearlyInterpolatedF0)
			e.numSamplesAfterPulseGenerated -= int(linearlyInterpolatedF0)
		} else {
			// absence
			excite[i] = 0.0
		}
		e.numSamplesAfterPulseGenerated++
	}

	return excite
}
