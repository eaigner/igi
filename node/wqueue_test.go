package node

import (
	"testing"
)

func TestWeightQueue(t *testing.T) {
	q := NewWeightQueue(10)

	for i := 1; i <= 10; i++ {
		if !q.Push(i, i) {
			t.Fatal()
		}
	}

	for i := 10; i > 0; i-- {
		i := q.Pop().(int)
		if i != i {
			t.Fatal(i)
		}
	}
}
