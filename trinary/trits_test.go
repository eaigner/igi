package trinary

import (
	"testing"
)

func TestTritsFromInt8(t *testing.T) {
	type test struct {
		in     []int8
		result bool
	}
	table := []test{
		{[]int8{-1, 0, 1}, true},
		{[]int8{-1, 0, 2}, false},
		{[]int8{-2, 0, 1}, false},
		{[]int8{-2}, false},
		{[]int8{-1}, true},
		{[]int8{}, true},
	}

	for _, v := range table {
		var trits Trits
		ok := TritsFromInt8(v.in, &trits)
		if ok != v.result {
			t.Fail()
		}
		if ok && trits.Len() != len(v.in) {
			t.Fail()
		}
		if !ok && trits.Len() != 0 {
			t.Fail()
		}
	}
}

func TestTritsAt(t *testing.T) {
	var trits Trits
	a := []int8{-1, 0, 1, 0, -1}
	if !TritsFromInt8(a, &trits) {
		t.Fail()
	}
	for i, v := range a {
		if trits.At(i) != v {
			t.Fail()
		}
	}
}

func TestTritsToTrytes(t *testing.T) {
	var in []int8

	r := []int8{-1, 0, 1}

	for _, i := range r {
		for _, j := range r {
			for _, k := range r {
				in = append(in, k, j, i)
			}
		}
	}

	if len(in) != len(tryteAlphabet)*3 {
		t.Fatal("invalid input length")
	}

	var v Trits
	if !TritsFromInt8(in, &v) {
		t.Fatal("could not convert to trits")
	}

	s := v.ToTrytes()
	expect := "NOPQRSTUVWXYZ9ABCDEFGHIJKLM"

	if s != expect {
		t.Logf("is=%s, want=%s", s, expect)
	}
}
