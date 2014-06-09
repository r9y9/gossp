package vocoder

import (
	"github.com/r9y9/gossp/excite"
	"github.com/r9y9/gossp/f0"
	"github.com/r9y9/gossp/io"
	"github.com/r9y9/gossp/mgcep"
	"github.com/r9y9/gossp/stft"
	"github.com/r9y9/gossp/window"
	"log"
	"testing"
)

func TestMGLSASynthesis(t *testing.T) {
	var (
		testData   []float64
		frameShift = 80
		frameLen   = 512
		alpha      = 0.41
		stage      = 12
		gamma      = -1.0 / float64(stage)
		order      = 24
		f0Seq      []float64
		ex         []float64
		mgc        [][]float64
	)

	w, err := io.ReadWav("../test_files/test16k.wav")
	if err != nil {
		log.Fatal(err)
	}
	testData = w.GetMonoData()

	// F0
	f0Seq = f0.SWIPE(testData, 16000, frameShift, 60.0, 700.0)

	// MGCep
	s := &stft.STFT{
		FrameShift: frameShift,
		FrameLen:   frameLen,
	}

	numFrames := s.NumFrames(testData)
	mgc = make([][]float64, numFrames)
	for i := 0; i < numFrames; i++ {
		windowed := window.BlackmanNormalized(s.FrameAtIndex(testData, i))
		mgc[i] = mgcep.MGCep(windowed, order, alpha, gamma)
	}

	// adjast number of frames
	m := min(len(f0Seq), len(mgc))
	f0Seq, mgc = f0Seq[:m], mgc[:m]

	// Excitation
	g := &excite.PulseExcite{
		SampleRate: 16000,
		FrameShift: frameShift,
	}
	ex = g.Generate(f0Seq)

	// Waveform generation
	synth := NewMGLSASpeechSynthesizer(order, alpha, stage, frameShift)

	_, err = synth.Synthesis(ex, mgc)
	if err != nil {
		t.Errorf("Error %s, want non-nil.", err)
	}
	// TODO(ryuichi) valid check
}
