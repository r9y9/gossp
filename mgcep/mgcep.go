// Package mgcep provides support for Mel-Generalized Cepstrum Analysis.
package mgcep

import (
	"github.com/mjibson/go-dsp/fft"
	"github.com/r9y9/gossp"
	"github.com/r9y9/gossp/3rdparty/mcepalpha"
	"github.com/r9y9/gossp/3rdparty/sptk"
)

// MGCep performs Mel-Generalized Cepstrum analysis and returns mgcep
// coefficients.
// Parameters:
//      order: order of mel-generalized cepstrum
//      alpha: frequency warping parameter
//      gamma: parameter of generalized log function
// Return:
//      mel-generalized cepstrum coefficients (length of order+1)
func MGCep(audioBuffer []float64,
	order int, alpha, gamma float64) []float64 {
	return sptk.MGCepWithDefaultParameters(audioBuffer, order, alpha, gamma)
}

// Calculate all-pass constant (alpha) for a given sample frequency.
func CalcMCepAlpha(sampleRate int) float64 {
	return mcepalpha.CalcMcepAlpha(sampleRate)
}

// Amplitude2Cepstrum returns (real) cepstrum from an amplitude spectrum.
func LogAmplitude2Cepsturm(logAmp []float64) []float64 {
	return gossp.ToReal(fft.IFFTReal(logAmp))
}

// Cepstrum2Amplitude return log amplitude spectrum from cepstrum.
// Note: This function requires symmetric cepstrum as the argument.
func Cepstrum2LogAmplitude(ceps []float64) []float64 {
	return gossp.ToReal(fft.FFTReal(ceps))
}

// LogAmp2MCep returns mel-cepstrum from a given log amplitude spectrum.
// Note: This function requires symmetric log-amplitude that length correspond
// to the number of frequency bins in fourier analysis.
func LogAmp2MCep(logAmp []float64, order int, alpha float64) []float64 {
	ceps := LogAmplitude2Cepsturm(logAmp)
	ceps[0] /= 2.0

	return FrequencyWarping(ceps, order, alpha)
}

// MCep2LogAmp returns log amplitude spectrum from a given mel-cepstrum.
func MCep2LogAmp(melCeps []float64, numFreqBins int, alpha float64) []float64 {
	// Convert to mel-scale from linear scale
	ceps := FrequencyWarping(melCeps, len(melCeps)-1, alpha)
	ceps[0] *= 2.0

	// Adjast number of frequeny bins
	actualCeps := make([]float64, numFreqBins)
	actualCeps[0] = ceps[0]
	for i := 1; i < len(ceps); i++ {
		actualCeps[i] = ceps[i]
		actualCeps[numFreqBins-i] = ceps[i]
	}

	return Cepstrum2LogAmplitude(actualCeps)
}
