package cwt

import (
	"math"
	"testing"
)

func TestCentScales(t *testing.T) {
	centInterval := 25
	minFreq, maxFreq := 55.0, 7040.0

	centScales := CreateCentScales(minFreq, maxFreq, centInterval)
	expected := 1.0 / 55.0
	if result := centScales[0]; result != expected {
		t.Errorf("centScales[0] where maxFreq %v, centInterval %v and minFreq %v = %v, want %v",
			maxFreq, centInterval, minFreq, result, expected)
	}

	expectedLogFreqDelta := float64(centInterval) / float64(CentPerOctave)
	for i := 1; i < len(centScales); i++ {
		delta := (math.Log(1.0/centScales[i]) - math.Log(1.0/centScales[i-1])) / math.Log(2.0)
		if math.Abs(delta-expectedLogFreqDelta) > 1.0e-10 {
			t.Errorf("Log delta of successive frequencies in cent-scale where maxFreq %v, centInterval %v and minFreq %v = %v, want %v",
				maxFreq, centInterval, minFreq, delta, expectedLogFreqDelta)
		}
	}
}

/*
func TestCentScaleIndexAndFreqConversion(t *testing.T) {
	// Create Gabor wavelet
	scale := &CentScale{
		BaseFreq:     55.0,
		NumOctaves:   7,
		CentInterval: 25,
	}

	// index to frequency conversion
	for index := 0; index < 5; index++ {
		actualFreq := scale.IndexToFreq(index)
		expectedFreq := scale.BaseFreq *
			math.Pow(2, float64(index)*float64(scale.CentInterval)/
				float64(CentPerOctave))

		if math.Abs(expectedFreq-actualFreq) > 3.0 {
			t.Errorf("Index %v corresponds to %f in Hz, want %f",
				index, expectedFreq, actualFreq)
		}
	}

	// frequency to index conversion
	NumBinsPerOctave := int(float64(CentPerOctave) / float64(scale.CentInterval))
	for multiplier := 1; multiplier < scale.NumOctaves-1; multiplier++ {
		freq := scale.BaseFreq * math.Pow(2, float64(multiplier-1))
		actualIndex := scale.FreqToIndex(freq)
		expectedIndex := NumBinsPerOctave * (multiplier - 1)
		if expectedIndex != actualIndex {
			t.Errorf("Frequency %f Hz corresponds to %v in index, want %v",
				freq, actualIndex, expectedIndex)
		}
	}
}
*/
