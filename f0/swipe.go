package f0

import (
	"github.com/r9y9/gossp/3rdparty/sptk"
)

// SWIPE: A SAWTOOTH WAVEFORM INSPIRED PITCH ESTIMATOR
// FOR SPEECH AND MUSIC.
// http://www.cise.ufl.edu/~acamacho/publications/dissertation.pdf
func SWIPE(audioBuffer []float64, sampleRate int, frameShift int,
	min, max float64) []float64 {
	return sptk.SWIPEWithDefaultParameters(audioBuffer,
		sampleRate,
		frameShift,
		min,
		max)
}
