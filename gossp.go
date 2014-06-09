// Package gossp provides support for speech signal processing.
package gossp

import (
	"math"
)

// DivideFrames returns overlapping divided frames.
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

// ToReal performs a real sequence from a complex sequence.
func ToReal(x []complex128) []float64 {
	y := make([]float64, len(x))
	for i, val := range x {
		y[i] = real(val)
	}
	return y
}

// Angle returns angle of complex number in radian.
func Angle(z complex128) float64 {
	return math.Atan2(imag(z), real(z))
}

// UnWrap returns unwraped phase given a phase spectrum.
func UnWrap(phase []float64) []float64 {
	unwraped := make([]float64, len(phase))

	unwraped[0] = phase[0]
	for i := 1; i < len(phase); i++ {
		diff := phase[i] - phase[i-1]
		switch {
		case diff > math.Pi:
			unwraped[i] = phase[i] - 2*math.Pi
		case diff < -math.Pi:
			unwraped[i] = phase[i] + 2*math.Pi
		default:
			unwraped[i] = phase[i]
		}
	}

	return unwraped
}

// Matlab
func HistC(x, edges []float64) []int {
	index := make([]int, len(edges))

	count := 1

	previousIndex := 0
	for i, val := range edges {
		previousIndex = i
		index[i] = 1
		if val >= x[0] {
			break
		}
	}

	for i := previousIndex; i < len(edges); i++ {
		if edges[i] < x[count] {
			index[i] = count
		} else {
			index[i] = count
			i--
			count++
		}

		previousIndex = i
		if count == len(x) {
			break
		}
	}
	count--

	for i := previousIndex + 1; i < len(edges); i++ {
		index[i] = count
	}

	return index
}

// Matlab
func Interp1(x, y, xi []float64) []float64 {
	yi := make([]float64, len(xi))

	h := make([]float64, len(x)-1)
	p := make([]float64, len(xi))
	s := make([]float64, len(xi))

	for i := 0; i < len(x)-1; i++ {
		h[i] = x[i+1] - x[i]
	}

	for i := 0; i < len(x); i++ {
		p[i] = float64(i)
	}

	k := HistC(x, xi)

	for i := range xi {
		s[i] = (xi[i] - x[k[i]-1]) / h[k[i]-1]
	}

	for i := range xi {
		yi[i] = y[k[i]-1] + s[i]*(y[k[i]]-y[k[i]-1])
	}

	return yi
}

// Symmetrize returns symmetrized vector given a input vector.
func Symmetrize(x []float64) []float64 {
	N := len(x)
	y := make([]float64, (N-1)*2)

	y[0] = x[0]
	for i := 1; i < N; i++ {
		y[i] = x[i]
		y[len(y)-i] = x[i]
	}

	return y
}

// Hz2Mel converts frequency in hz to mel.
func Hz2Mel(freq float64) float64 {
	return 1127.01048 * math.Log(freq/700.0+1.0)
}

// Mel2Hz converts frequency in mel to hz.
func Mel2Hz(mel float64) float64 {
	return 700.0 * (math.Exp(mel/1127.01048) - 1.0)
}
