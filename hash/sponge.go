package hash

type Sponge interface {
	Absorb(v []int8)
	Squeeze(v []int8)
	Reset(mode int)
}
