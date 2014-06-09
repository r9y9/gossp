// Package stft provides support for Short-Time Fourier Transform (STFT)
// Analysis.
package stft

import (
	"github.com/mjibson/go-dsp/fft"
	"github.com/r9y9/gossp/window"
)

// Reference:
// D. W. Griffin and J. S. Lim, "Signal estimation from modified short-time
// Fourier transform," IEEE Trans. ASSP, vol.32, no.2, pp.236–243, Apr. 1984.

// STFT represents Short Time Fourier Transform Analysis.
type STFT struct {
	FrameShift int
	FrameLen   int
	Window     []float64 // window funtion
}

// New returns a new STFT instance.
func New(frameShift, frameLen int) *STFT {
	s := &STFT{
		FrameShift: frameShift,
		FrameLen:   frameLen,
		Window:     window.CreateHanning(frameLen),
	}

	return s
}

// NumFrames returnrs the number of frames that will be analyzed in STFT.
func (s *STFT) NumFrames(input []float64) int {
	return int(float64(len(input)-s.FrameLen)/float64(s.FrameShift)) + 1
}

// DivideFrames returns overlapping divided frames for STFT.
func (s *STFT) DivideFrames(input []float64) [][]float64 {
	numFrames := s.NumFrames(input)
	frames := make([][]float64, numFrames)
	for i := 0; i < numFrames; i++ {
		frames[i] = s.FrameAtIndex(input, i)
	}
	return frames
}

// FrameAtIndex returns frame at specified index given an input signal.
// Note that it doesn't make copy of input.
func (s *STFT) FrameAtIndex(input []float64, index int) []float64 {
	return input[index*s.FrameShift : index*s.FrameShift+s.FrameLen]
}

// STFT returns complex spectrogram given an input signal.
func (s *STFT) STFT(input []float64) [][]complex128 {
	numFrames := s.NumFrames(input)
	spectrogram := make([][]complex128, numFrames)

	frames := s.DivideFrames(input)
	for i, frame := range frames {
		// Windowing
		windowed := window.Windowing(frame, s.Window)
		// Complex Spectrum
		spectrogram[i] = fft.FFTReal(windowed)
	}

	return spectrogram
}

// ISTFT performs invere STFT signal reconstruction and returns reconstructed
// signal.
func (s *STFT) ISTFT(spectrogram [][]complex128) []float64 {
	frameLen := len(spectrogram[0])
	numFrames := len(spectrogram)
	reconstructedSignal := make([]float64, frameLen+numFrames*s.FrameShift)

	// Griffin's method
	windowSum := make([]float64, len(reconstructedSignal))
	for i := 0; i < numFrames; i++ {
		buf := fft.IFFT(spectrogram[i])
		index := 0
		for t := i * s.FrameShift; t < i*s.FrameShift+frameLen; t++ {
			reconstructedSignal[t] += real(buf[index]) * s.Window[index]
			windowSum[t] += s.Window[index] * s.Window[index]
			index++
		}
	}

	// Normalize by window
	for n := range reconstructedSignal {
		reconstructedSignal[n] /= windowSum[n]
	}

	return reconstructedSignal
}
