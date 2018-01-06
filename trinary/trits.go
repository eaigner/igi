package trinary

import (
	"sync"
)

const (
	tryteAlphabet = "9ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type Trits struct {
	a   []int8
	s   string
	mtx sync.Mutex
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

func (t *Trits) ToTrytes() string {
	t.mtx.Lock()
	defer func() {
		t.mtx.Unlock()
	}()
	if len(t.s) > 0 {
		return t.s
	}
	if t.Len()%3 != 0 {
		return ""
	}
	n := t.Len() / 3
	v := make([]byte, n)
	for i := 0; i < n; i++ {
		j := t.a[i*3] + t.a[i*3+1]*3 + t.a[i*3+2]*9
		if j < 0 {
			j += int8(len(tryteAlphabet))
		}
		v[i] = tryteAlphabet[j]
	}
	t.s = string(v)
	return t.s
}
