package vocoder

// B2MC performs conversion from MLSA filter coefficients to Mel-cepstrum.
func B2MC(b []float64, alpha float64) []float64 {
	mc := make([]float64, len(b))
	m := len(b) - 1

	mc[m] = b[m]
	d := mc[m]

	for i := m - 1; i >= 0; i-- {
		o := b[i] + alpha*d
		d = b[i]
		mc[i] = o
	}

	return mc
}
