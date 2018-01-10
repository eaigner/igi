package hash

import "testing"

func TestWeightMagnitude(t *testing.T) {
	if w := WeightMagnitude([]int8{1 - 1, 1, 0}); w != 1 {
		t.Fatal(w)
	}
	if w := WeightMagnitude([]int8{1 - 1, 1, 0, 0, 0, 0}); w != 4 {
		t.Fatal(w)
	}
	if w := WeightMagnitude([]int8{1 - 1, 1}); w != 0 {
		t.Fatal(w)
	}
	if w := WeightMagnitude([]int8{}); w != 0 {
		t.Fatal(w)
	}
}
