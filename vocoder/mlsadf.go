package vocoder

// MLSAFilter represents a Mel-Log Spectrum Approximation (MLSA) filter,
// which is composed of two cascade filters.
type MLSAFilter struct {
	baseFilter1 *MLSACascadeFilter
	baseFilter2 *MLSACascadeFilter
}

// MLSACascadeFilter represents a cascade filter which contains MLSA base filters.
type MLSACascadeFilter struct {
	cascadeFilters []*MLSABaseFilter
	padeCoef       []float64
	delay          []float64
}

// MLSABaseFilter represents a base filter of the MLSA digital filter.
type MLSABaseFilter struct {
	Order int // Order of mel-cepstrum
	// all-pass constant (alpha) that determines how frequency axis is warped.
	Alpha float64

	delay []float64
}

// NewMLSABaseFilter returns its instance.
// It requires the order of mel-cepstrum and all-pass constant (alpha).
func NewMLSABaseFilter(order int, alpha float64) *MLSABaseFilter {
	mlsaBase := new(MLSABaseFilter)
	mlsaBase.Order = order
	mlsaBase.Alpha = alpha
	mlsaBase.delay = make([]float64, order+1)
	return mlsaBase
}

// Filter returns filtered sample given the inpuy sample and filter coeffficients.
func (bf *MLSABaseFilter) Filter(sample float64, filterCoef []float64) float64 {
	bf.delay[0] = sample
	bf.delay[1] = (1.0-bf.Alpha*bf.Alpha)*bf.delay[0] + bf.Alpha*bf.delay[1]

	result := 0.0
	for i := 2; i < len(filterCoef); i++ {
		bf.delay[i] = bf.delay[i] + bf.Alpha*(bf.delay[i+1]-bf.delay[i-1])
		result += bf.delay[i] * filterCoef[i]
	}

	// Special case
	if len(filterCoef) == 2 {
		result += bf.delay[1] * filterCoef[1]
	}

	// t <- t+1 in time
	for i := len(bf.delay) - 1; i > 1; i-- {
		bf.delay[i] = bf.delay[i-1]
	}

	return result
}

// NewMLSABaseFilter returns its instance and error with
// the order of mel-cepstrum, all-pass constant (alpha) and the order
// of pade approximation.
// Order of pade approximation 4 or 5 is only supported.
// If other one is specified, it returns non-nil error.
func NewMLSACascadeFilter(order int, alpha float64, orderOfPade int) *MLSACascadeFilter {
	cf := new(MLSACascadeFilter)
	cf.cascadeFilters = make([]*MLSABaseFilter, orderOfPade+1)
	cf.delay = make([]float64, orderOfPade+1)
	cf.padeCoef = make([]float64, orderOfPade+1)

	for i := 0; i <= orderOfPade; i++ {
		cf.cascadeFilters[i] = NewMLSABaseFilter(order, alpha)
	}

	switch orderOfPade {
	case 4:
		cf.padeCoef[0] = 1.0
		cf.padeCoef[1] = 4.999273e-1
		cf.padeCoef[2] = 1.067005e-1
		cf.padeCoef[3] = 1.170221e-2
		cf.padeCoef[4] = 5.656279e-4
	case 5:
		cf.padeCoef[0] = 1.0
		cf.padeCoef[1] = 4.999391e-1
		cf.padeCoef[2] = 1.107098e-1
		cf.padeCoef[3] = 1.369984e-2
		cf.padeCoef[4] = 9.564853e-4
		cf.padeCoef[5] = 3.041721e-5
	default:
		panic("MLSA digital filter: Order of pade approximation must be 4 or 5")
	}

	return cf
}

// Filter returns filtered sample given the inpuy sample and filter coeffficients.
func (cf *MLSACascadeFilter) Filter(sample float64, filterCoef []float64) float64 {
	result, feedBack := 0.0, 0.0

	for i := len(cf.padeCoef) - 1; i >= 1; i-- {
		cf.delay[i] = cf.cascadeFilters[i].Filter(cf.delay[i-1], filterCoef)
		val := cf.delay[i] * cf.padeCoef[i]
		if i%2 == 1 {
			feedBack += val
		} else {
			feedBack -= val
		}
		result += val
	}

	cf.delay[0] = feedBack + sample
	result += cf.delay[0]

	return result
}

// NewMLSAFilter returns its instance and error.
// It requires the order of mel-cepstrum, all-pass constant (alpha)
// and the order of pade approximation.
func NewMLSAFilter(order int, alpha float64, orderOfPade int) *MLSAFilter {
	mlsa := new(MLSAFilter)

	// First stage filter
	mlsa.baseFilter1 = NewMLSACascadeFilter(2, alpha, orderOfPade)

	// Second stage filter
	mlsa.baseFilter2 = NewMLSACascadeFilter(order+1, alpha, orderOfPade)

	return mlsa
}

// Filter returns filtered sample given the input and MLSA filter coefficients.
// The MLSA filter is composed of two stage cascade filter.
func (f *MLSAFilter) Filter(sample float64, filterCoef []float64) float64 {
	firstStageCoef := []float64{0, filterCoef[1]}

	// two stage cascade filtering
	return f.baseFilter2.Filter(f.baseFilter1.Filter(sample, firstStageCoef), filterCoef)
}
