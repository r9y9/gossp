package vocoder

import (
	"github.com/r9y9/gossp/mgcep"
	"math"
)

// MGCep2MGLSAFilterCoef performs a conversion from mel-generalized cepstrum
// to MGLSA filter coefficients.
func MGCep2MGLSAFilterCoef(melgcep []float64, alpha, gamma float64) []float64 {
	filterCoef := MCep2MLSAFilterCoef(melgcep, alpha)
	filterCoef = mgcep.GNorm(filterCoef, gamma)

	// scale by gamma
	filterCoef[0] = math.Log(filterCoef[0])
	for i := 1; i < len(filterCoef); i++ {
		filterCoef[i] *= gamma
	}

	return filterCoef
}
