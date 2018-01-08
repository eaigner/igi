package trinary

const (
	tryteAlphabet = "9ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	tritsPerByte  = 5
	tritsPerTryte = 3
	tritRadix     = 3
	maxTritValue  = (tritRadix - 1) / 2
	minTritValue  = -maxTritValue
)

var (
	bytesToTrits [243][tritsPerByte]int8
)

func init() {
	var trits [tritsPerByte]int8
	for i := 0; i < 243; i++ {
		copy(bytesToTrits[i][:], trits[:])
		incrementTrits(trits[:], tritsPerByte)
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
func Trits(dst []int8, src []byte) int {
	n := len(src) * tritsPerByte
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
			return 0
		}

		copy(dst[offset:offset+o], bytesToTrits[x][0:o])

		offset += tritsPerByte
	}

	for offset < n {
		dst[offset] = 0
		offset++
	}

	return n
}

func validTrit(v int8) bool {
	return v >= -1 && v <= 1
}

// LenBytes returns the number of bytes that trits encodes to.
func LenBytes(trits []int8) int {
	return (len(trits) + tritsPerByte - 1) / tritsPerByte
}

// LenTrits returns the number of trits that fit into bytes.
func LenTrits(bytes []byte) int {
	return len(bytes) * tritsPerByte
}

// Bytes converts a trit to a byte slice.
// dst must at least be LenBytes(src) long.
// Returns the number of bytes written.
func Bytes(dst []byte, src []int8) int {
	n := LenBytes(src)

	if len(dst) < n {
		return 0
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
			v = v*int8(tritRadix) + src[i*tritsPerByte+j]
		}
		dst[i] = byte(v)
	}

	return n
}

// Trytes converts trits into a tryte string.
func Trytes(t []int8) string {
	if len(t)%3 != 0 {
		return ""
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
	return string(v)
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
