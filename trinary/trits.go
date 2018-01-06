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
