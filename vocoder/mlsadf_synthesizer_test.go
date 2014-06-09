package vocoder

import (
	"github.com/r9y9/gossp/excite"
	"github.com/r9y9/gossp/f0"
	"github.com/r9y9/gossp/io"
	"github.com/r9y9/gossp/mgcep"
	"github.com/r9y9/gossp/stft"
	"github.com/r9y9/gossp/window"
	"log"
	"math"
	"testing"
)

func TestMLSASynthesis(t *testing.T) {
	var (
		testData   []float64
		frameShift = 80
		frameLen   = 512
		alpha      = 0.41
		order      = 24
		pd         = 5
		f0Seq      []float64
		ex         []float64
		mc         [][]float64
	)

	w, err := io.ReadWav("../test_files/test16k.wav")
	if err != nil {
		log.Fatal(err)
	}
	testData = w.GetMonoData()

	// F0
	f0Seq = f0.SWIPE(testData, 16000, frameShift, 60.0, 700.0)

	// Mcep
	s := &stft.STFT{
		FrameShift: frameShift,
		FrameLen:   frameLen}

	numFrames := s.NumFrames(testData)
	mc = make([][]float64, numFrames)
	for i := 0; i < numFrames; i++ {
		windowed := window.BlackmanNormalized(s.FrameAtIndex(testData, i))
		mc[i] = mgcep.MCep(windowed, order, alpha)
	}

	// adjast number of frames
	m := min(len(f0Seq), len(mc))
	f0Seq, mc = f0Seq[:m], mc[:m]

	// Excitation
	g := &excite.PulseExcite{
		SampleRate: 16000,
		FrameShift: frameShift,
	}
	ex = g.Generate(f0Seq)

	// Waveform generation
	synth := NewMLSASpeechSynthesizer(order, alpha, pd, frameShift)

	_, err = synth.Synthesis(ex, mc)
	if err != nil {
		t.Errorf("Error %s, want non-nil.", err)
	}
	// TODO(ryuichi) valid check
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func createSin(freq float64, sampleRate, length int) []float64 {
	sin := make([]float64, length)
	for i := range sin {
		sin[i] = math.Sin(2.0 * math.Pi * freq * float64(i) / float64(sampleRate))
	}
	return sin
}
