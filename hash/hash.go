package hash

import "bytes"

const (
	SizeBytes = 49 // a hash would only require 46 bytes, but curl needs more to compute
	SizeTrits = 243
)

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

// Valid returns true if the provided byte hash has the correct length and is not the zero hash.
func Valid(hash []byte) bool {
	return len(hash) == SizeBytes && !Zero(hash)
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

// Valid returns true if the provided int8/trit hash has the correct length and is not the zero hash.
func ValidInt8(hash []int8) bool {
	return len(hash) == SizeTrits && !ZeroInt8(hash)
}

// ToBytes converts an int8 hash to bytes
func ToBytes(hash []int8) []byte {
	buf := make([]byte, len(hash))
	for i, v := range hash {
		buf[i] = byte(v)
	}
	return buf
}

// ToInt8 converts a byte hash to int8
func ToInt8(hash []byte) []int8 {
	buf := make([]int8, len(hash))
	for i, v := range hash {
		buf[i] = int8(v)
	}
	return buf
}
