package hash

const (
	CurlP27 = 27
	CurlP81 = 81
)

const (
	curlStateLength = 3 * SizeTrits
)

var (
	truthTable = [...]int{1, 0, -1, 2, 1, -1, 0, 2, -1, 1, 0}
)

type Curl struct {
	mode    int
	state   [curlStateLength]int
	scratch [curlStateLength]int
}

func (c *Curl) Reset(mode int) {
	c.mode = mode
	memsetRepeat(c.state[:], 0)
	memsetRepeat(c.scratch[:], 0)
}

func (c *Curl) Absorb(v []int8) {
	in := v
	var n int
	for len(in) > 0 {
		n = len(in)
		if n > SizeTrits {
			n = SizeTrits
		}
		for i, v := range in[:n] {
			c.state[i] = int(v)
		}
		c.transform()
		in = in[n:]
	}
}

func (c *Curl) Squeeze(v []int8) {
	in := v
	var n int
	for len(in) > 0 {
		n = len(in)
		if n > SizeTrits {
			n = SizeTrits
		}
		for i, v := range c.state[:n] {
			in[i] = int8(v)
		}
		c.transform()
		in = in[n:]
	}
}

func (c *Curl) transform() {
	var i, j int
	for round := 0; round < c.mode; round++ {
		c.scratch = c.state // copy
		for k := 0; k < curlStateLength; k++ {
			j = i
			if i < 365 {
				i += 364
			} else {
				i -= 365
			}
			c.state[k] = truthTable[c.scratch[j]+(c.scratch[i]<<2)+5]
		}
	}
}

func memsetRepeat(a []int, v int) {
	if len(a) == 0 {
		return
	}
	a[0] = v
	for bp := 1; bp < len(a); bp *= 2 {
		copy(a[bp:], a[:bp])
	}
}
