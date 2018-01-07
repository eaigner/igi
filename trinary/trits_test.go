package trinary

import (
	"testing"
)

func TestValidate(t *testing.T) {
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
		if Validate(v.in) != v.result {
			t.Fatal()
		}
	}
}

func TestTrytes(t *testing.T) {
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

	s := Trytes(in)
	expect := "NOPQRSTUVWXYZ9ABCDEFGHIJKLM"

	if s != expect {
		t.Fatalf("is=%s, want=%s", s, expect)
	}
}

func TestBytes(t *testing.T) {
	// Must be a multiple of tritsPerByte
	var a = []int8{-1, 0, 1, 0, -1, 1, 1, 1, 1, -1}
	var b = make([]int8, 10)

	if x := LenBytes(a); x != 2 {
		t.Fatal(x)
	}

	var buf [10]byte

	if x := Bytes(buf[:], a); x != 2 {
		t.Fatal(x)
	}
	if x := Trits(b, buf[:2]); x != len(a) {
		t.Fatal(x)
	}
	if !Equals(a, b) {
		t.Fatal()
	}
}

func TestEquals(t *testing.T) {
	// TODO(era): impl
}
