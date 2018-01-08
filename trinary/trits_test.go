package trinary

import (
	"crypto/rand"
	"testing"
)

var (
	trits10 = []int8{-1, 0, 1, 0, -1, 1, 1, 1, 1, -1}
	bytes10 = []byte{0xb7, 0xd7}
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
	expect := "NOPQRSTUVWXYZ9ABCDEFGHIJKLM"

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

	s, err := Trytes(in)

	if err != nil {
		t.Fatal(err)
	}
	if s != expect {
		t.Fatalf("is=%s, want=%s", s, expect)
	}
}

func TestTritsSliceTooSmall(t *testing.T) {
	var dst []int8
	n, err := Trits(dst, bytes10)

	if err != errBufferTooSmall {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatal(n)
	}
	if len(dst) != 0 {
		t.Fatal(dst)
	}
}

func TestTritsSliceValid(t *testing.T) {
	max := LenTrits(bytes10)

	if max != 10 {
		t.Fatal(max)
	}

	dst := make([]int8, max)
	n, err := Trits(dst, bytes10)

	if err != nil {
		t.Fatal(err)
	}
	if n != 10 {
		t.Fatal(n)
	}
	if !Equals(dst, trits10) {
		t.Fatal(dst)
	}
}

func TestTritsMaliciousBytes(t *testing.T) {
	var buf = make([]byte, 10)
	var max = LenTrits(buf)
	var tbuf = make([]int8, max)

	for i := 0; i < 1000; i++ {
		n, err := rand.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		if n != len(buf) {
			t.Fatal(n)
		}

		nt, err := Trits(tbuf, buf)

		if nt == 0 && err != nil {
			// ok
		} else if nt == max && err == nil {
			// ok
		} else {
			t.Fatal(nt, err)
		}
	}
}

func TestBytes(t *testing.T) {
	// Must be a multiple of tritsPerByte
	var b = make([]int8, len(trits10))

	if x := LenBytes(trits10); x != 2 {
		t.Fatal(x)
	}

	var buf [10]byte

	x, err := Bytes(buf[:], trits10)

	if err != nil {
		t.Fatal(err)
	}
	if x != 2 {
		t.Fatal(x)
	}

	x, err = Trits(b, buf[:2])

	if err != nil {
		t.Fatal(err)
	}
	if x != len(trits10) {
		t.Fatal(x)
	}
	if !Equals(trits10, b) {
		t.Fatal()
	}
}

func TestEquals(t *testing.T) {
	if !Equals(trits10, trits10) {
		t.Fatal()
	}
	if Equals(trits10, trits10[:2]) {
		t.Fatal()
	}
}
