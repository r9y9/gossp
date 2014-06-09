package vocoder

import (
	"errors"
	"math"
)

// MGLSASpeechSynthesizer represents a speech synthesizer based on the
// MGLSA Filter.
type MGLSASpeechSynthesizer struct {
	FrameShift int
	Alpha      float64 // all-pass constant
	Gamma      float64 // parameter of generalized logarithmic function
	NumStage   int
	coreFilter *MGLSAFilter // used in sample by sample waveform generation
}

// NewMGLSASpeechSynthesizer returns its instance given parameters.
func NewMGLSASpeechSynthesizer(order int, alpha float64, numStage int,
	frameShift int) *MGLSASpeechSynthesizer {
	synthesizer := &MGLSASpeechSynthesizer{
		FrameShift: frameShift,
		Alpha:      alpha,
		NumStage:   numStage,
		Gamma:      -1.0 / float64(numStage),
		coreFilter: NewMGLSAFilter(order, alpha, numStage),
	}

	return synthesizer
}

// Synthesis synthesizes a speech signal from an excitation signal and
// corresponding mel-ceptrum sequence.
func (s *MGLSASpeechSynthesizer) Synthesis(excite []float64,
	mgcepSequence [][]float64) ([]float64, error) {
	if len(excite) != len(mgcepSequence)*s.FrameShift {
		return nil, errors.New("MGLSA Speech Synthesis: The length of excitation and mel-generalized cepstrum times frame peroid doesn't match.")
	}

	// synthesized speech signal will be stored
	synthesizedSpeech := make([]float64, len(excite))

	previousMgcep := mgcepSequence[0]
	for i, currentMgcep := range mgcepSequence {
		if i > 0 {
			previousMgcep = mgcepSequence[i-1]
		}
		// Synthesize a part of speech
		startIndex, endIndex := i*s.FrameShift, (i+1)*s.FrameShift
		partOfSpeech := s.SynthesisOneFrame(excite[startIndex:endIndex], previousMgcep, currentMgcep)

		for j, val := range partOfSpeech {
			synthesizedSpeech[i*s.FrameShift+j] = val
		}
	}

	return synthesizedSpeech, nil
}

// SynthesisOneFrame synthesizes a part of speech signal from an excitation signal
// and succesive two mel-cepstrum sequence. It requires all-pass constant (alpha).
// Mel-cepstral coefficients between two succesive mel-cepstrum are linearly
// interpolated.
func (s *MGLSASpeechSynthesizer) SynthesisOneFrame(excite []float64,
	previousMgcep, currentMgcep []float64) []float64 {
	// Convert to MGLSA filter coefficients from Mel-cepstrum
	currentFilterCoef := MGCep2MGLSAFilterCoef(currentMgcep, s.Alpha, s.Gamma)
	previousFilterCoef := MGCep2MGLSAFilterCoef(previousMgcep, s.Alpha, s.Gamma)

	// Compute slope
	slope := make([]float64, len(currentMgcep))
	for i := 0; i < len(slope); i++ {
		slope[i] = (currentFilterCoef[i] - previousFilterCoef[i]) / float64(len(excite))
	}

	partOfSpeech := make([]float64, len(excite))
	linearlyInterpolatedCoef := make([]float64, len(previousFilterCoef))
	copy(linearlyInterpolatedCoef, previousFilterCoef)
	for i := 0; i < len(excite); i++ {
		// Multyply power coeffcient
		scaledExcitation := excite[i] * math.Exp(linearlyInterpolatedCoef[0])
		// Filtering
		partOfSpeech[i] = s.coreFilter.Filter(scaledExcitation, linearlyInterpolatedCoef)
		// Linear interpolation of filter coefficients
		for j := 0; j < len(slope); j++ {
			linearlyInterpolatedCoef[j] += slope[j]
		}
	}

	return partOfSpeech
}
