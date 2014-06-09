package mgcep

// GC2GC peforms Generalized cepstrum transformation.
func GC2GC(c1 []float64, gamma1 float64, m2 int, gamma2 float64) []float64 {
	m1 := len(c1)
	c2 := make([]float64, m2+1)

	c2[0] = c1[0]

	for m := 1; m <= m2; m++ {
		ss1, ss2 := 0.0, 0.0
		min := m1
		if m1 >= m {
			min = m - 1
		}
		for k := 1; k <= min; k++ {
			cc := c1[k] * c2[m-k]
			ss2 += float64(k) * cc
			ss1 += float64(m-k) * cc
		}

		if m <= m1 {
			c2[m] = c1[m] + (gamma2*ss2-gamma1*ss1)/float64(m)
		} else {
			c2[m] = (gamma2*ss2 - gamma1*ss1) / float64(m)
		}
	}

	return c2
}
