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

type Trits struct {
	a []int8
}

func TritsFromInt8(a []int8, t *Trits) bool {
	for _, v := range a {
		if !validTrit(v) {
			return false
		}
	}
	copy(t.a, a)
	return true
}

func BytesToTrits(b []byte, t *Trits) int {
	// const trinarySize = 8019
	n := len(b) * tritsPerByte
	offset := 0

	if len(t.a) < n {
		t.a = make([]int8, n)
	}

	for i := 0; i < len(b) && offset < n; i++ {
		x := int(int8(b[i])) // cast twice or we will lose the sign
		if x < 0 {
			x += len(bytesToTrits)
		}
		o := tritsPerByte
		j := n - offset
		if j < tritsPerByte {
			o = j
		}

		// TODO: handle case of malicious bytes greater than len(bytesToTrits)

		copy(t.a[offset:offset+o], bytesToTrits[x][0:o])

		offset += tritsPerByte
	}

	for offset < n {
		t.a[offset] = 0
		offset++
	}

	return n
}

func validTrit(v int8) bool {
	return v >= -1 && v <= 1
}

func (t Trits) Len() int {
	return len(t.a)
}

func (t Trits) ToTrytes() string {
	if len(t.a)%3 != 0 {
		return ""
	}
	n := len(t.a) / 3
	v := make([]byte, n)
	for i := 0; i < n; i++ {
		j := t.a[i*3] + t.a[i*3+1]*3 + t.a[i*3+2]*9
		if j < 0 {
			j += int8(len(tryteAlphabet))
		}
		v[i] = tryteAlphabet[j]
	}
	return string(v)
}
