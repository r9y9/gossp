package mgcep

import (
	"github.com/r9y9/gossp/3rdparty/sptk"
)

func MCep(audioBuffer []float64, order int, alpha float64) []float64 {
	return sptk.MCepWithDefaultParameters(audioBuffer, order, alpha)
}

func MCep2Energy(mc []float64, alpha float64, length int) float64 {
	// Back to linear cepsrum
	c := FreqT(mc, length-1, -alpha)

	impulseResponse := C2IR(c, length)

	energy := 0.0
	for k := range impulseResponse {
		energy += impulseResponse[k] * impulseResponse[k]
	}

	return energy
}
