package trinary

import (
	"reflect"
	"testing"
)

func TestToTrits(t *testing.T) {
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
		trits, ok := ToTrits(v.in)
		if ok != v.result {
			t.Fail()
		}
		if ok && !reflect.DeepEqual([]int8(trits), v.in) {
			t.Fail()
		}
		if !ok && len(trits) > 0 {
			t.Fail()
		}
	}
}
