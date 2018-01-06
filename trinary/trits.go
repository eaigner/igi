package trinary

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
