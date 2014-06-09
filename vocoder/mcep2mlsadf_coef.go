package vocoder

// MCep2MLSAFilterCoef performs a conversion from mel-cepstrum to MLSA filter
// coefficients given mel-cepstrum and all-pass constant(alpha).
func MCep2MLSAFilterCoef(mcep []float64, alpha float64) []float64 {
	filterCoef := make([]float64, len(mcep))

	filterCoef[len(filterCoef)-1] = mcep[len(filterCoef)-1]

	for i := len(filterCoef) - 2; i >= 0; i-- {
		filterCoef[i] = mcep[i] - alpha*filterCoef[i+1]
	}

	return filterCoef
}

// Alias to SPTK
// will be removed
func MC2B(mcep []float64, alpha float64) []float64 {
	return MCep2MLSAFilterCoef(mcep, alpha)
}
