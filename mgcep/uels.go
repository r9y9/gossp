package mgcep

import (
	"github.com/r9y9/gossp/3rdparty/sptk"
)

// TODO(ryuichi) replace with pure Go.
func UELS(audioBuffer []float64, order int) []float64 {
	return sptk.UELSSimple(audioBuffer, order)
}
