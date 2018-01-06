package trinary

type Trits []int8

func ToTrits(v []int8) (Trits, bool) {
	for _, v := range v {
		if !validTrit(v) {
			return []int8{}, false
		}
	}
	return Trits(v), true
}

func validTrit(v int8) bool {
	return v >= -1 && v <= 1
}
