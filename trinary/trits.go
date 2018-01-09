package trinary

import "errors"

const (
	tryteAlphabet      = "9ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	tritsPerByte       = 5
	tritsPerTryte      = 3
	tritRadix     int8 = 3
	maxTritValue  int8 = (tritRadix - 1) / 2
	minTritValue  int8 = -maxTritValue
)

var (
	tryteRuneIndex = make(map[rune]int, len(tryteAlphabet))
	bytesToTrits   [243][tritsPerByte]int8
	trytesToTrits  [27][tritsPerTryte]int8
)

var (
	errBufferTooSmall = errors.New("buffer too small")
	errMultipleThree  = errors.New("buffer size must be a multiple of 3")
)

func init() {
	var trits [tritsPerByte]int8
	for i := 0; i < 243; i++ {
		copy(bytesToTrits[i][:], trits[:tritsPerByte])
		incrementTrits(trits[:], tritsPerByte)
	}
	for i := 0; i < 27; i++ {
		copy(trytesToTrits[i][:], trits[:tritsPerTryte])
		incrementTrits(trits[:], tritsPerTryte)
	}
	for i, c := range tryteAlphabet {
		tryteRuneIndex[c] = i
	}
}

func incrementTrits(trits []int8, n int) {
	for i := 0; i < n; i++ {
		trits[i]++
		if trits[i] > maxTritValue {
			trits[i] = minTritValue
		} else {
			break
		}
	}
}

// Validate checks if a trit slice is valid.
func Validate(a []int8) bool {
	for _, v := range a {
		if !validTrit(v) {
			return false
		}
	}
	return true
}

// Trits converts a byte to a trit slice.
// dst must at least be LenTrits(src) long.
// Returns the number of trits written.
func Trits(dst []int8, src []byte) (int, error) {
	n := len(src) * tritsPerByte

	if n < len(dst) {
		return 0, errBufferTooSmall
	}

	offset := 0

	for i := 0; i < len(src) && offset < n; i++ {
		x := int(int8(src[i])) // cast twice or we will lose the sign
		if x < 0 {
			x += len(bytesToTrits)
		}
		o := tritsPerByte
		j := n - offset
		if j < tritsPerByte {
			o = j
		}

		// check if we have space left
		if len(dst) < offset+o {
			return 0, errBufferTooSmall
		}

		copy(dst[offset:offset+o], bytesToTrits[x][0:o])

		offset += tritsPerByte
	}

	for offset < n {
		dst[offset] = 0
		offset++
	}

	return n, nil
}

// TritsFromTrytes converts a tryte string to trits.
// dst must at least be LenTritsFromTrytes(src) long.
// Returns the number of trits written
func TritsFromTrytes(dst []int8, src string) (int, error) {
	n := LenTritsFromTrytes(len(src))
	if len(dst) < n {
		return 0, errBufferTooSmall
	}
	for i, c := range src {
		j := i * tritsPerTryte
		copy(dst[j:j+tritsPerTryte], trytesToTrits[tryteRuneIndex[c]][:])
	}
	return n, nil
}

func validTrit(v int8) bool {
	return v >= -1 && v <= 1
}

// LenBytes returns the number of bytes for n trits
func LenBytes(n int) int {
	return (n + tritsPerByte - 1) / tritsPerByte
}

// LenTrits returns the number of trits for n bytes
func LenTrits(n int) int {
	return n * tritsPerByte
}

// LenTritsFromTrytes returns the number of trits for n trytes
func LenTritsFromTrytes(n int) int {
	return n * tritsPerTryte
}

// Bytes converts a trit to a byte slice.
// dst must at least be LenBytes(src) long.
// Returns the number of bytes written.
func Bytes(dst []byte, src []int8) (int, error) {
	n := LenBytes(len(src))

	if len(dst) < n {
		return 0, errBufferTooSmall
	}

	var v int8
	var j0 int

	for i := 0; i < n; i++ {
		v = 0
		j0 = len(src) - i*tritsPerByte
		if j0 >= 5 {
			j0 = tritsPerByte
		}
		for j := j0 - 1; j >= 0; j-- {
			v = v*tritRadix + src[i*tritsPerByte+j]
		}
		dst[i] = byte(v)
	}

	return n, nil
}

// Trytes converts trits into a tryte string.
func Trytes(t []int8) (string, error) {
	if len(t)%3 != 0 {
		return "", errMultipleThree
	}
	n := len(t) / 3
	v := make([]byte, n)
	for i := 0; i < n; i++ {
		j := t[i*3] + t[i*3+1]*3 + t[i*3+2]*9
		if j < 0 {
			j += int8(len(tryteAlphabet))
		}
		v[i] = tryteAlphabet[j]
	}
	return string(v), nil
}

// Int64 converts trits to an int64 value.
func Int64(t []int8) int64 {
	var v int64 = 0
	for i := len(t) - 1; i >= 0; i-- {
		v = v*int64(tritRadix) + int64(t[i])
	}
	return v
}

// Equals compares two trit buffers.
func Equals(t1 []int8, t2 []int8) bool {
	if len(t1) != len(t2) {
		return false
	}
	for i, v := range t1 {
		if t2[i] != v {
			return false
		}
	}
	return true
}
