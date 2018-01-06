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
	type test struct {
		in     []int8
		result string
	}

	var table []test

	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			for k := -1; k < 2; k++ {
				v := []int8{int8(k), int8(j), int8(i)}
				c := tryteAlphabet[(len(table)+14)%len(tryteAlphabet)]
				table = append(table, test{v, string(c)})
			}
		}
	}

	if len(table) != len(tryteAlphabet) {
		t.Fatal("invalid table length")
	}

	for _, v := range table {
		var trits Trits
		if !TritsFromInt8(v.in, &trits) {
			t.Fail()
		}
		var trytes Trytes
		if !trits.ToTrytes(&trytes) {
			t.Fail()
		}
		if s := trytes.String(); s != v.result {
			t.Logf("is=%s, want=%s", s, v.result)
		}
	}
}
