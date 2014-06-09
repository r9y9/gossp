package mgcep

// Alias
func FreqT(ceps []float64, order int, alpha float64) []float64 {
	return FrequencyWarping(ceps, order, alpha)
}

// FrequencyWarping returns frequency-warped (mel) cesptrum given a cepstrum.
// The order is a desired order of the cepstrum and alhpa is an all-pass
// constant.
func FrequencyWarping(ceps []float64, order int, alpha float64) []float64 {
	warpedCeps := make([]float64, order+1)
	prev := make([]float64, order+1)

	m1 := len(ceps) - 1

	for i := -m1; i <= 0; i++ {
		copy(prev, warpedCeps)
		if order >= 0 {
			warpedCeps[0] = ceps[-i] + alpha*prev[0]
		}
		if order >= 1 {
			warpedCeps[1] = (1.0-alpha*alpha)*prev[0] + alpha*prev[1]
		}
		for m := 2; m <= order; m++ {
			warpedCeps[m] = prev[m-1] + alpha*(prev[m]-warpedCeps[m-1])
		}
	}

	return warpedCeps
}
