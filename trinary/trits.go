package trinary

const (
	tryteAlphabet = "9ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type Trits struct {
	a []int8
}

func TritsFromInt8(a []int8, t *Trits) bool {
	for _, v := range a {
		if !validTrit(v) {
			return false
		}
	}
	t.a = a
	return true
}

func validTrit(v int8) bool {
	return v >= -1 && v <= 1
}

func (t Trits) Len() int {
	return len(t.a)
}

func (t Trits) At(i int) int8 {
	return t.a[i]
}

func (t Trits) ToTrytes(v *Trytes) bool {
	if t.Len()%3 != 0 {
		return false
	}
	n := t.Len() / 3
	v.buf = make([]byte, n)
	for i := 0; i < n; i++ {
		j := t.At(i*3) + t.At(i*3+1)*3 + t.At(i*3+2)*9
		if j < 0 {
			j += int8(len(tryteAlphabet))
		}
		v.buf[i] = tryteAlphabet[j]
	}
	return true
}

type Trytes struct {
	buf []byte
}

func TrytesFromInt8(a []int8, t *Trytes) bool {
	var v Trits
	return TritsFromInt8(a, &v) && v.ToTrytes(t)
}

func (t Trytes) String() string {
	return string(t.buf)
}
