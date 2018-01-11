package hash

import (
	"errors"

	"github.com/eaigner/igi/trinary"
)

var (
	errBufferTooSmall = errors.New("buffer too small")
)

type Sponge interface {
	Absorb(v []int8)
	Squeeze(v []int8)
	Reset(mode int)
}

// SqueezeBytes calls Squeeze and writes the result as bytes to b.
// b must be SizeBytes long.
func SqueezeBytes(s Sponge, b []byte) (int, error) {
	if len(b) < SizeBytes {
		return 0, errBufferTooSmall
	}
	var t [SizeTrits]int8

	s.Squeeze(t[:])

	n, err := trinary.Bytes(b, t[:])
	if err != nil {
		return 0, err
	}
	return n, nil
}
