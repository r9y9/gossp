package vocoder

// MGLSAFilter represents a Mel-generalized log-spectral approximation filter.
type MGLSAFilter struct {
	cascadeFilters []*MGLSABaseFilter
	numStage       int // -1.0/gamma
}

// MGLSABaseFilter represents a base filter of MGLSAFilter.
type MGLSABaseFilter struct {
	Order int // Order of mel-cepstrum
	// all-pass constant (alpha) that determines how frequency axis is warped.
	Alpha float64

	delay []float64
}

func NewMGLSABaseFilter(order int, alpha float64) *MGLSABaseFilter {
	f := &MGLSABaseFilter{
		Order: order,
		Alpha: alpha,
		delay: make([]float64, order+1),
	}
	return f
}

func (bf *MGLSABaseFilter) Filter(sample float64, filterCoef []float64) float64 {
	y := bf.delay[0] * filterCoef[1]

	for i := 1; i < len(filterCoef)-1; i++ {
		bf.delay[i] += bf.Alpha * (bf.delay[i+1] - bf.delay[i-1])
		y += bf.delay[i] * filterCoef[i+1]
	}

	result := sample - y

	// t <- t+1 in time
	for i := len(bf.delay) - 1; i > 0; i-- {
		bf.delay[i] = bf.delay[i-1]
	}
	bf.delay[0] = bf.Alpha*bf.delay[0] + (1.0-bf.Alpha*bf.Alpha)*result

	return result
}

func NewMGLSAFilter(order int, alpha float64, numStage int) *MGLSAFilter {
	f := &MGLSAFilter{
		numStage: numStage,
	}

	for i := 0; i < numStage; i++ {
		base := NewMGLSABaseFilter(order, alpha)
		f.cascadeFilters = append(f.cascadeFilters, base)
	}

	return f
}

func (mf *MGLSAFilter) Filter(sample float64, filterCoef []float64) float64 {
	// Initial value
	filterd := sample

	for i := 0; i < mf.numStage; i++ {
		filterd = mf.cascadeFilters[i].Filter(filterd, filterCoef)
	}

	return filterd
}
