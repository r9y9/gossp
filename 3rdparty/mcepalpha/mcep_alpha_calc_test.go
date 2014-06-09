package mcepalpha

// Writtern by Ryuichi Yamamoto

import (
	"math"
	"testing"
)

func testBase(t *testing.T, sampleRate int, expected float64) {
	actualValue := CalcMcepAlpha(sampleRate)
	if math.Abs(expected-actualValue) > 1.0e-10 {
		t.Errorf("All-pass constant alpha for the sample frequency %f is %f, want %f.",
			sampleRate, actualValue, expected)
	}
}

func TestMcepAlpha(t *testing.T) {
	testBase(t, 8000, 0.312)
	testBase(t, 11025, 0.357)
	testBase(t, 16000, 0.41)
	testBase(t, 22050, 0.455)
	testBase(t, 32000, 0.504)
	testBase(t, 44100, 0.544)
	testBase(t, 48000, 0.554)
}
