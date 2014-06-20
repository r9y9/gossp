package stft

import (
	"github.com/r9y9/gossp/io"
	"github.com/r9y9/gossp/window"
	"log"
	"math"
	"testing"
)

func loadReal16kData() []float64 {
	w, err := io.ReadWav("../test_files/test16k.wav")
	if err != nil {
		log.Fatal(err)
	}
	return w.GetMonoData()
}

func prepareWindowFunctions(frameLen int) [][]float64 {
	windowFunctions := make([][]float64, 4)
	windowFunctions[0] = window.CreateHanning(frameLen)
	windowFunctions[1] = window.CreateHamming(frameLen)
	windowFunctions[2] = window.CreateBlackman(frameLen)
	windowFunctions[3] = window.CreateGaussian(frameLen, 0.4)
	return windowFunctions
}

func TestConsistencyBetweenSTFTAndISTFT(t *testing.T) {
	var (
		testData            = loadReal16kData()
		testFrameLen        = []int{4096, 2048, 1024, 512}
		testFrameShiftDenom = []int{2, 3, 4, 5, 6, 7, 8} // 50% overlap 75% ...
		errTolerance        = 1.2
	)

	for _, frameLen := range testFrameLen {
		windowFunctions := prepareWindowFunctions(frameLen)
		for _, win := range windowFunctions {
			for _, denom := range testFrameShiftDenom {
				frameShift := frameLen / denom
				s := &STFT{
					FrameShift: frameShift,
					FrameLen:   frameLen,
					Window:     win,
				}
				reconstructed := s.ISTFT(s.STFT(testData))
				if containNAN(reconstructed) {
					t.Errorf("NAN contained, want non NAN contained.")
				}

				err := absErr(reconstructed, testData)
				if err > errTolerance {
					t.Errorf("[Frame length %d, Frame shift %d] %f error, want less than %f", frameLen, frameShift, err, errTolerance)
				}
			}
		}
	}
}

func containNAN(a []float64) bool {
	for _, val := range a {
		if math.IsNaN(val) {
			return true
		}
	}
	return false
}

func absErr(a, b []float64) float64 {
	length := len(a)
	if len(b) < length {
		length = len(b)
	}
	err := 0.0
	for i := 0; i < length; i++ {
		err += math.Abs(a[i] - b[i])
	}
	return err / float64(length)
}
