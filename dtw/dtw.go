// Package dtw provides support for Dynamic Time Warping (DTW) to align two
// time-series.
package dtw

// DTW represents dynamic time warping that can be used to align two
// time-series.
type DTW struct {
	ForwardStep      int
	BackwardStep     int
	Template         [][]float64
	CostTable        [][]float64
	BackTracePointer [][]int
}

// SetTemplate sets the template to align an input sequence.
func (d *DTW) SetTemplate(template [][]float64) {
	d.Template = template
	d.CostTable = make([][]float64, 1)
	d.BackTracePointer = make([][]int, 1)

	h := len(template)
	d.CostTable[0] = make([]float64, h)
	d.BackTracePointer[0] = make([]int, h)

	for i := 0; i < h; i++ {
		d.CostTable[0][i] = float64(i)
		d.BackTracePointer[0][i] = i
	}
}

// DTW peforms dynamic programming to align an input sequence to template.
func (d *DTW) DTW(sequence [][]float64) []int {
	// Forward recursion
	for i := 0; i < len(sequence); i++ {
		d.Update(sequence[i])
	}

	return d.BackTrace()
}

// TODO(ryuichi) use interface
func (d *DTW) ObservationCost(v []float64, index int) float64 {
	t := d.Template[index]
	cost := 0.0
	for i := range t {
		cost += (v[i] - t[i]) * (v[i] - t[i])
	}
	return cost
}

// TODO(ryuichi) use interface
func (d *DTW) TransitionCost(i, j int) float64 {
	switch {
	case j == i+1:
		return 0
	case i == j:
		return 1
	default:
		return 2
	}
}

// Update performs one stop of forward recursion.
func (d *DTW) Update(v []float64) {
	h := len(d.CostTable[0])
	lastCost := d.CostTable[len(d.CostTable)-1]
	currentCost := make([]float64, len(lastCost))
	currentBackTracePointer := make([]int, len(lastCost))

	for i := 0; i < h; i++ {
		minIndex := i
		obs, trans := d.ObservationCost(v, i), d.TransitionCost(minIndex, i)
		minCost := lastCost[minIndex] + obs + trans

		for j := i - d.BackwardStep; j <= i+d.ForwardStep; j++ {
			if j < 0 || j > h-1 {
				continue
			}
			cost := lastCost[j] + obs + d.TransitionCost(j, i)
			if cost < minCost {
				minCost, minIndex = cost, j
			}
		}
		currentCost[i] = minCost
		currentBackTracePointer[i] = minIndex
	}

	d.CostTable = append(d.CostTable, currentCost)
	d.BackTracePointer = append(d.BackTracePointer, currentBackTracePointer)
}

func (d *DTW) findMinIndex() int {
	lastCost := d.CostTable[len(d.CostTable)-1]
	minIndex, min := 0, lastCost[0]
	for i, val := range lastCost {
		if val <= min {
			minIndex, min = i, val
		}
	}
	return minIndex
}

// BackTrace finds a minumum cost path.
func (d *DTW) BackTrace() []int {
	timeLength := len(d.CostTable) - 1 // minus 1 (initial state)
	minPath := make([]int, timeLength)

	// Find the last index of min-cost path
	minPath[timeLength-1] = d.findMinIndex()

	// Recursive backtrace
	for i := timeLength - 1; i > 0; i-- {
		minPath[i-1] = d.BackTracePointer[i+1][minPath[i]]
	}

	return minPath
}
