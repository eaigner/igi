package hash

import "bytes"

var (
	nullBytes = make([]byte, SizeBytes)
)

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

// Zero returns true if hash contains the zero hash.
func Zero(hash []byte) bool {
	return bytes.Equal(hash, nullBytes)
}

// ZeroInt8 returns true if hash contains the zero trit hash.
func ZeroInt8(hash []int8) bool {
	if len(hash) != SizeTrits {
		return false
	}
	for _, v := range hash {
		if v != 0 {
			return false
		}
	}
	return true
}
