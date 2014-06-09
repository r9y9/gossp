package cwt

// TODO(ryuichi) scale and frequency conversion

import (
	"math"
)

const (
	CentPerSemitone = 100                         // cent per semiton
	OctaveDiv       = 12                          // num octaves
	CentPerOctave   = OctaveDiv * CentPerSemitone // 1200
)

// WaveletScaler represents scales that can be used in Wavelet analysis.
type WaveletScaler interface {
	Scales() []float64
	DeltaFreq() float64
}

// CentScale represents cent-based (logarighmically spaced) scales
// for Wavelet analysis.
type CentScale struct {
	MinFreq, MaxFreq float64
	CentInterval     int
	scales           []float64 // scale = 1/frequency
}

type WaveletScale struct {
	S0     float64
	Dj     float64
	J      int
	scales []float64
}

func (w *WaveletScale) Scales() []float64 {
	if w.scales == nil {
		w.scales = CreateWaveletScales(w.S0, w.Dj, w.J)
	}
	return w.scales
}

func (w *WaveletScale) DeltaFreq() float64 {
	return w.Dj
}

func CreateWaveletScales(s0, dj float64, J int) []float64 {
	scales := make([]float64, J)
	for j := 0; j < len(scales); j++ {
		scales[j] = s0 * math.Pow(2, float64(j)*dj)
	}
	return scales
}

// NewCentScale returns a new cent scale with default parameters.
func NewCentScale() *CentScale {
	c := &CentScale{
		MinFreq:      27.5,
		MaxFreq:      5000.0,
		CentInterval: 25, // 4 bins per octave
	}
	return c
}

// Scales returns a set of scales.
func (c *CentScale) Scales() []float64 {
	if c.scales == nil {
		c.scales = CreateCentScales(c.MinFreq, c.MaxFreq, c.CentInterval)
	}
	return c.scales
}

// DeltaFreq returns delta frequency.
func (c *CentScale) DeltaFreq() float64 {
	return float64(c.CentInterval) / float64(CentPerOctave)
}

// FreqToIndex converts frequency [Hz] to the corresponding index.
func (c *CentScale) FreqToIndex(freq float64) int {
	return int(float64(CentPerOctave) /
		float64(c.CentInterval) * math.Log2(freq/c.MinFreq))
}

// IndexToFreq converts index to the corresponding frequency [Hz]
func (c *CentScale) IndexToFreq(index int) float64 {
	return c.MinFreq * math.Pow(2.0, float64(index)*
		float64(c.CentInterval)/float64(CentPerOctave))
}

// CreateCentScales creates cent-based scales for wavelet analysis.
func CreateCentScales(minFreq float64, maxFreq float64,
	centInterval int) []float64 {
	numFreqBins := int(CentPerOctave*(math.Log(maxFreq/minFreq))/
		float64(centInterval)/math.Log(2.0)) + 1
	scales := make([]float64, numFreqBins)

	for s := 0; s < numFreqBins; s++ {
		order := float64(s) * float64(centInterval) / float64(CentPerOctave)
		freq := minFreq * math.Pow(2.0, order)
		scales[s] = 1.0 / freq
	}

	return scales
}
