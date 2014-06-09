package mgcep

// MGC2MGC performs Mel generalized cepstrum transformation.
func MGC2MGC(c1 []float64, alpha1, gamma1 float64,
	m2 int, alpha2, gamma2 float64) []float64 {
	c2 := make([]float64, m2+1)

	alpha := (alpha2 - alpha1) / (1.0 - alpha1*alpha2)

	if alpha == 0.0 {
		c2 = GNorm(c1, gamma1)
		c2 = GC2GC(c2, gamma1, m2, gamma2)
		c2 = IGNorm(c2, gamma2)
		return c2
	} else {
		c2 = FrequencyWarping(c1, m2+1, alpha)
		c2 = GNorm(c2, gamma1)
		c2 = GC2GC(c2, gamma1, m2, gamma2)
		c2 = IGNorm(c2, gamma2)
		return c2
	}
}
