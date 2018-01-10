package hash

// WeightMagnitude returns the weight magnitude of trits.
func WeightMagnitude(trits []int8) int {
	last := len(trits) - 1
	for i := last; i > 0; i-- {
		if trits[i] != 0 {
			return last - i
		}
	}
	return 0
}
