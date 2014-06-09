package cwt

import (
	"math"
)

// PointGaborWavelet represents Gabor Wavelet Transform that is focued on
// a specific point.
type PointGaborWavelet struct {
	SampleRate int
	FrameShift int
	Std        float64 // standard deviation of gaussian window
	Scale      WaveletScaler
	MinFreq    float64 // Used for computing expected window length

	// tables to avoid dupulicate operations
	preComputed bool
	gaborWaveletTables
}

type gaborWaveletTables struct {
	gauss [][]float64
	cos   [][]float64
	sin   [][]float64
}

func NewPointGaborWavelet(sampleRate, frameShift int) *PointGaborWavelet {
	scale := NewCentScale()

	return &PointGaborWavelet{
		SampleRate: sampleRate,
		FrameShift: frameShift,
		Std:        3.0,
		Scale:      scale,
		MinFreq:    scale.MinFreq,
	}
}

// CWT performs Continuous Wavelet Transform (CWT) given an input signal
// with a fucus on the specified position (second argument).
// This function returns CWT complex spectrum
func (g *PointGaborWavelet) CWT(input []float64, center int) []complex128 {
	g.checkBeforeCWT()

	numFreqBins := g.NumFreqBins()
	result := make([]complex128, numFreqBins)
	scales := g.Scale.Scales()

	// convolution in time-domain
	for y := 0; y < numFreqBins; y++ {
		// window lengh [sec]
		dt := scales[y] * g.Std * math.Sqrt(-2.0*math.Log(0.01))
		// window length [sample]
		dx := int(dt * float64(g.SampleRate))
		real_wt, imag_wt := 0.0, 0.0
		for m := -dx; m <= dx; m++ {
			n := center + m
			if n >= 0 && n < len(input) {
				real_wt += input[n] * g.gauss[y][m+dx] * g.cos[y][m+dx]
				imag_wt += input[n] * g.gauss[y][m+dx] * g.sin[y][m+dx]
			}
		}
		norm := 1.0 / math.Sqrt(scales[y])
		result[y] = complex(norm*real_wt, norm*imag_wt)
	}

	return result
}

// NumFreqBins returns the number of frequency bins that will be computed
// in CWT analysis.
func (g *PointGaborWavelet) NumFreqBins() int {
	return len(g.Scale.Scales())
}

// NumFrames returns the number of frames in CWT time-frequency analysis
func (g *PointGaborWavelet) NumFrames(input []float64) int {
	return int(float64(len(input))/float64(g.FrameShift)) + 1
}

func (g *PointGaborWavelet) checkBeforeCWT() {
	if !g.preComputed {
		g.createTables()
		g.preComputed = true
	}
}

// CreateTables computes gaussian, sin and cos values that will be used in CWT.
func (g *PointGaborWavelet) createTables() {
	// expected number of frequency bins
	numFreqBins := g.NumFreqBins()
	scales := g.Scale.Scales()

	g.gauss = make([][]float64, numFreqBins)
	g.cos = make([][]float64, numFreqBins)
	g.sin = make([][]float64, numFreqBins)

	// expected window length
	winLen := int(2.0*float64(g.SampleRate)*1.0/g.MinFreq*g.Std*
		math.Sqrt(-2.0*math.Log(0.01)) + 1.0)

	for y := 0; y < numFreqBins; y++ {
		g.gauss[y] = make([]float64, winLen)
		g.cos[y] = make([]float64, winLen)
		g.sin[y] = make([]float64, winLen)
	}

	norm := 1.0 / math.Sqrt(2.0*math.Pi*g.Std*g.Std)
	for y := 0; y < numFreqBins; y++ {
		dt := scales[y] * g.Std * math.Sqrt(-2.0*math.Log(0.01))
		dx := int(dt * float64(g.SampleRate))
		for m := -dx; m <= dx; m++ {
			t := float64(m) / float64(g.SampleRate) / scales[y]
			g.gauss[y][m+dx] = norm *
				math.Exp(-t*t/(2.0*g.Std*g.Std))
			g.sin[y][m+dx] = math.Sin(2.0 * math.Pi * t)
			g.cos[y][m+dx] = math.Cos(2.0 * math.Pi * t)
		}
	}
}
