package vocoder

import (
	"errors"
	"math"
)

// MLSASpeechSynthesizer represents a speech synthesizer based on
// MLSA Filter.
type MLSASpeechSynthesizer struct {
	CoreFilter *MLSAFilter // used in sample by sample speech generation
	FrameShift int
	Alpha      float64 // all-pass constant
}

// NewMLSASpeechSynthesizer returns its instance given parameters.
func NewMLSASpeechSynthesizer(numMceps int, alpha float64, orderOfPade int,
	frameShift int) *MLSASpeechSynthesizer {
	synthesizer := new(MLSASpeechSynthesizer)

	synthesizer.CoreFilter = NewMLSAFilter(numMceps, alpha, orderOfPade)
	synthesizer.FrameShift = frameShift
	synthesizer.Alpha = alpha

	return synthesizer
}

// Synthesis synthesizes a speech signal from an excitation signal and
// corresponding mel-ceptrum sequence.
func (s *MLSASpeechSynthesizer) Synthesis(excite []float64, mcepSequence [][]float64) ([]float64, error) {
	if len(excite) != len(mcepSequence)*s.FrameShift {
		return nil, errors.New("MLSA Speech Synthesis: length of excitation and mel-cepstrum times frame peroid doesn't match.")
	}

	// synthesized speech signal will be stored
	synthesizedSpeech := make([]float64, len(excite))

	previousMcep := mcepSequence[0]
	for i, currentMcep := range mcepSequence {
		if i > 0 {
			previousMcep = mcepSequence[i-1]
		}
		// Synthesize a part of speech
		startIndex, endIndex := i*s.FrameShift, (i+1)*s.FrameShift
		partOfSpeech := s.SynthesisOneFrame(excite[startIndex:endIndex], previousMcep, currentMcep)

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
func (s *MLSASpeechSynthesizer) SynthesisOneFrame(excite []float64,
	previousMcep, currentMcep []float64) []float64 {
	// Convert to MLSA filter coefficients from Mel-cepstrum
	currentFilterCoef := MCep2MLSAFilterCoef(currentMcep, s.Alpha)
	previousFilterCoef := MCep2MLSAFilterCoef(previousMcep, s.Alpha)

	// Compute slope
	slope := make([]float64, len(currentMcep))
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
		partOfSpeech[i] = s.CoreFilter.Filter(scaledExcitation, linearlyInterpolatedCoef)
		// Linear interpolation of filter coefficients
		for j := 0; j < len(slope); j++ {
			linearlyInterpolatedCoef[j] += slope[j]
		}
	}

	return partOfSpeech
}
