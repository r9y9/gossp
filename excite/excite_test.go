package excite

import (
	"math"
	"testing"
)

func TestLength(t *testing.T) {
	// setup
	sampleRate := 44100
	frameShift := 441
	ex := &PulseExcite{
		SampleRate: sampleRate,
		FrameShift: frameShift,
	}

	// create dummy data
	f0Sequence := []float64{100, 100, 99, 98, 97, 0, 0, 0, 0, 0, 0, 70, 70, 71, 72, 72, 73, 75, 75}

	// Generate excitation
	excitation := ex.Generate(f0Sequence)

	// check
	expectedLen := frameShift * len(f0Sequence)
	if len(excitation) != expectedLen {
		t.Errorf("The length of generated excitaiton is %d, want %d", len(excitation), expectedLen)
	}
}

func TestExcitationPeriodicity(t *testing.T) {
	// setup
	sampleRate := 44100
	frameShift := 441
	ex := &PulseExcite{
		SampleRate: sampleRate,
		FrameShift: frameShift,
	}

	// create dummy data (100hz)
	f0Sequence := []float64{100.0, 100.0, 100.0, 100.0, 100.0}

	// ex should generate one pulse per 441(44100/100) samples
	excitation := ex.Generate(f0Sequence)

	// Check pulse value for the ideal pulse location
	normalizedF0 := float64(sampleRate) / 100.0
	expectedPulseValue := math.Sqrt(normalizedF0)
	if excitation[frameShift+1] != expectedPulseValue {
		t.Errorf("Generated pulse is %f, want %f.", excitation[frameShift], expectedPulseValue)
	}

	// Check periodicity
	numDetectedPulse := 1
	prevIndex := frameShift + 1
	for i := frameShift + 2; i < len(excitation); i++ {
		if excitation[i] == expectedPulseValue {
			if (i - prevIndex) != int(normalizedF0) {
				t.Errorf("Pulse should be generated per %d samples, but %d.", int(normalizedF0), i-prevIndex)
			}
			numDetectedPulse++
			prevIndex = i
		}
	}

	if numDetectedPulse != 4 {
		t.Errorf("4 Pulses should be generated, but %d.", numDetectedPulse)
	}
}

func TestConsistencyBetweenGaussianOrNot(t *testing.T) {
	// setup
	sampleRate := 44100
	frameShift := 441
	ex := &PulseExcite{
		SampleRate: sampleRate,
		FrameShift: frameShift,
	}

	// create dummy data
	f0Sequence := []float64{100.0, 100.0, 99.0, 98.0, 97.0, 0, 0, 0, 0, 0, 0, 70.0, 70.0, 71.0, 72.0, 72.0, 73.0, 75.0, 75.0}

	// Pseudo random samples for non-voice segments
	excitation1 := ex.Generate(f0Sequence)

	ex.UseGauss = true

	// Gaussian random samples for non-voice segments
	excitation2 := ex.Generate(f0Sequence)

	// check consistency
	for i := range excitation1 {
		if (i >= 0*frameShift && i < 5*frameShift) || i >= 12*frameShift {
			// check for voice-segments
			if excitation1[i] != excitation2[i] {
				t.Errorf("Voice segments: Index %d, Generated excitation must be equal regardless of type of randam sample (Gaussian or not). %f != %f", i, excitation1[i], excitation2[i])
			}
		}
	}
}
