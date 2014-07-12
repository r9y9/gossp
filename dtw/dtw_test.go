package dtw

import (
	"testing"
)

func TestDTW(t *testing.T) {
	v1 := [][]float64{{1, 2, 3}, {1, 2, 4}, {1, 8, 5}, {10, 3, 6}}
	v2 := [][]float64{{1, 2, 3}, {1, 2, 4}, {1, 2, 5}, {1, 8, 5}, {10, 3, 6}}

	d := &DTW{ForwardStep: 0, BackwardStep: 1}
	d.SetTemplate(v1)
	path := d.DTW(v2)

	expected := []int{0, 1, 1, 2, 3}
	for i := range expected {
		if expected[i] != path[i] {
			t.Error("Result", path, "Expected", expected)
			break
		}
	}

	z1 := [][]float64{{0}, {1}, {2}, {3}, {4}, {5}}
	z2 := [][]float64{{0}, {0}, {1}, {2}, {3}, {4}, {4}, {5}}

	d2 := &DTW{ForwardStep: 0, BackwardStep: 1}
	d2.SetTemplate(z1)
	path2 := d2.DTW(z2)

	expected2 := []int{0, 0, 1, 2, 3, 4, 4, 5}
	for i := range expected2 {
		if expected2[i] != path2[i] {
			t.Error("Result", path2, "Expected", expected2)
			break
		}
	}
}
