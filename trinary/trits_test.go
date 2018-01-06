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
