// Package timefreq provides support for time-frequency analysis.
package timefreq

import (
	"math"
	"math/cmplx"
)

// DivideFrames returns overlapping divided frames for STFT.
func DivideFrames(input []float64, frameLen, frameShift int) [][]float64 {
	numFrames := int(float64(len(input)-frameLen)/float64(frameShift)) + 1
	frames := make([][]float64, numFrames)
	for i := 0; i < numFrames; i++ {
		frames[i] = input[i*frameShift : i*frameShift+frameLen]
	}
	return frames
}

// SplitSpectrum splits complex spectrum X(k) to amplitude |X(k)|
// and angle(X(k))
func SplitSpectrum(spec []complex128) ([]float64, []float64) {
	amp := make([]float64, len(spec))
	phase := make([]float64, len(spec))
	for i, val := range spec {
		amp[i] = cmplx.Abs(val)
		phase[i] = math.Atan2(imag(val), real(val))
	}

	return amp, phase
}

// SplitSpectrogram returns SpilitSpectrum for each time frame.
func SplitSpectrogram(spectrogram [][]complex128) ([][]float64, [][]float64) {
	numFrames, numFreqBins := len(spectrogram), len(spectrogram[0])
	amp := create2DSlice(numFrames, numFreqBins)
	phase := create2DSlice(numFrames, numFreqBins)

	for i := 0; i < numFrames; i++ {
		amp[i], phase[i] = SplitSpectrum(spectrogram[i])
	}

	return amp, phase
}

func create2DSlice(rows, cols int) [][]float64 {
	s := make([][]float64, rows)
	for i := range s {
		s[i] = make([]float64, cols)
	}
	return s
}

// ReconstructSpectrum returns complex spectrum from amplitude
// and phase spectrum.
// angle(X(k)) and |X(k)| -> X(k)
func ReconstructSpectrum(amp, phase []float64) []complex128 {
	spec := make([]complex128, len(amp))
	for i := range amp {
		spec[i] = complex(amp[i], 0.0) * cmplx.Exp(complex(0.0, phase[i]))
	}
	return spec
}

// ReconstructSpectrogram returns complex spectrogram from amplitude
// phase spectrogram.
func ReconstructSpectrogram(amplitudeSpectrogram,
	phaseSpectrogram [][]float64) [][]complex128 {
	spectrogram := make([][]complex128, len(amplitudeSpectrogram))
	for i := range amplitudeSpectrogram {
		a, p := amplitudeSpectrogram[i], phaseSpectrogram[i]
		spectrogram[i] = ReconstructSpectrum(a, p)
	}
	return spectrogram
}
