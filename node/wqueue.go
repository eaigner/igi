package node

import (
	"container/heap"
	"sync"
)

type WeightQueue struct {
	q   pQueue
	c   chan bool
	mtx sync.Mutex
}

func NewWeightQueue(maxLen int) *WeightQueue {
	return &WeightQueue{
		q: make(pQueue, 0),
		c: make(chan bool, maxLen),
	}
}

func (q *WeightQueue) Push(value interface{}, weight int) bool {
	q.mtx.Lock()
	heap.Push(&q.q, &pqItem{
		value:    value,
		priority: weight,
	})
	q.mtx.Unlock()
	q.c <- true
	return true
}

func (q *WeightQueue) Pop() interface{} {
	<-q.c
	q.mtx.Lock()
	item := heap.Pop(&q.q).(*pqItem)
	q.mtx.Unlock()
	return item.value
}

type pqItem struct {
	value    interface{}
	priority int
	index    int
}

type pQueue []*pqItem

func (pq pQueue) Len() int { return len(pq) }

func (pq pQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq pQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *pQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*pqItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *pQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}
