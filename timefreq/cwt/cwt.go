// Package cwt provides support for Continuous Wavelet Transform (CWT).
package cwt

import (
	"github.com/mjibson/go-dsp/fft"
	"math"
)

// FrequencyReper represents interface that has frequency representation.
type FrequencyReper interface {
	FrequencyRep(omega []float64, scale float64) []float64
}

// AtTimeZeroer represents interface that can be evaluated at zero.
type AtTimeZeroer interface {
	AtTimeZero() float64
}

// ReconstructionFactorer represents interface that has reconstruction factor
// which is used  for Inverse Continuous Wavelet transform.
type ReconstructionFactorer interface {
	ReconstructionFactor() (float64, error)
}

// FourierPerioder represents interface that converts a wavelet scale to
// fourier period.
type FourierPerioder interface {
	FourierPeriod(scale float64) float64
}

// WaveletBasis represents interface that wavelet basis should satisfy.
type WaveletBasiser interface {
	FrequencyReper
	AtTimeZeroer
	ReconstructionFactorer
	FourierPerioder
}

// Wavelet represents Continuous Wavelet Analysis.
type Wavelet struct {
	DeltaT float64
	Scale  WaveletScaler
	Basis  WaveletBasiser
}

func NewWavelet(deltaT float64,
	basis WaveletBasiser, scale WaveletScaler) *Wavelet {
	w := &Wavelet{
		DeltaT: deltaT,
		Basis:  basis,
		Scale:  scale,
	}
	return w
}

func nextHighestPowOf2(n int) int {
	return int(math.Pow(2.0, math.Ceil(math.Log2(float64(n)))))
}

// CWT performs Continuous Wavelet Transform.
func (w *Wavelet) CWT(x []float64) [][]complex128 {
	scales := w.Scale.Scales()

	// Memoly allocation
	N := len(x)
	spectrogram := make([][]complex128, N)
	for i := range spectrogram {
		spectrogram[i] = make([]complex128, len(scales))
	}

	// Zero padding
	pN := nextHighestPowOf2(N)
	zeroPadded := make([]float64, pN)
	copy(zeroPadded, x)

	// FFT of input signal
	y := fft.FFTReal(zeroPadded)
	omega := angularFrequency(pN, w.DeltaT)

	// Convolution in Freuency domain
	for j, scale := range scales {
		norm := math.Sqrt(2.0 * math.Pi * scale / w.DeltaT)

		waveletInFreqDomain := w.Basis.FrequencyRep(omega, scale)

		product := make([]complex128, pN)
		for n := 0; n < pN; n++ {
			product[n] = complex(norm*waveletInFreqDomain[n], 0) * y[n]
		}
		convolved := fft.IFFT(product)

		for n := 0; n < N; n++ {
			spectrogram[n][j] = convolved[n]
		}
	}

	return spectrogram
}

// angularFrequency returns angular frequency used in Continuous Wavelet Analysis.
func angularFrequency(N int, dt float64) []float64 {
	omega := make([]float64, N)

	for i := 0; i <= N/2; i++ {
		omega[i] = 2.0 * math.Pi * float64(i) / (dt * float64(N))
	}
	for i := N/2 + 1; i < N; i++ {
		omega[i] = -2.0 * math.Pi * float64(i) / (dt * float64(N))
	}

	return omega
}

// ComputeReconstructionFactor returns reconstruction connstant for the
// Inverse CWT.
func (w *Wavelet) ComputeReconstructionFactor(N int) float64 {
	dj := w.Scale.DeltaFreq()
	omega := angularFrequency(N, w.DeltaT)

	scales := w.Scale.Scales()
	deltaCWT := make([]float64, len(scales))

	// Wavelet transform for delta function
	for j, scale := range scales {
		norm := math.Sqrt(2.0 * math.Pi * scale / w.DeltaT)
		w := w.Basis.FrequencyRep(omega, scale)
		for k := range w {
			deltaCWT[j] += w[k]
		}
		deltaCWT[j] *= norm / float64(N)
	}

	// Scale normalization
	C := 0.0
	for j := range deltaCWT {
		C += deltaCWT[j] / math.Sqrt(scales[j])
	}
	C *= dj * math.Sqrt(w.DeltaT) / w.Basis.AtTimeZero()

	return C
}

// ICWT performs Inverse Continous Wavelet Transform.
func (w *Wavelet) ICWT(spectrogram [][]complex128) []float64 {
	C, err := w.Basis.ReconstructionFactor()
	if err != nil {
		// Not defined in mother wavelet, so we can compute.
		C = w.ComputeReconstructionFactor(len(spectrogram))
	}

	dj := w.Scale.DeltaFreq()
	scales := w.Scale.Scales()

	norm := dj * math.Sqrt(w.DeltaT) / (C * w.Basis.AtTimeZero())

	x := make([]float64, len(spectrogram))
	for n, coef := range spectrogram {
		for j := range coef {
			x[n] += real(spectrogram[n][j]) / math.Sqrt(scales[j])
		}
		x[n] *= norm
	}

	return x
}

// FrouierFreq returns Fourier frequencies converted from wavelet scales.
func (w *Wavelet) FourierFreq() []float64 {
	scales := w.Scale.Scales()

	freq := make([]float64, len(scales))
	for i, scale := range scales {
		freq[i] = 1.0 / w.Basis.FourierPeriod(scale)
	}

	return freq
}
