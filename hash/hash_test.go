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

func TestZero(t *testing.T) {
	if Zero([]byte{}) {
		t.Fatal()
	}
	if Zero([]byte{1, 0, 1}) {
		t.Fatal()
	}

	b := make([]byte, SizeBytes)

	if !Zero(b) {
		t.Fatal()
	}

	b[3] = 1

	if Zero(b) {
		t.Fatal()
	}
}

func TestZeroInt8(t *testing.T) {
	if ZeroInt8([]int8{}) {
		t.Fatal()
	}
	if ZeroInt8([]int8{1, 0, -1}) {
		t.Fatal()
	}

	b := make([]int8, SizeTrits)

	if !ZeroInt8(b) {
		t.Fatal()
	}

	b[3] = 1

	if ZeroInt8(b) {
		t.Fatal()
	}
}
