package queue

import (
	"container/heap"
	"sync"
)

// WeightQueue implements a weighted priority queue.
// In a high load regime, the network connection would be saturated and transactions with higher MWM
// (minimum weight magnitude) would rise to the top of the priority queue, so if you want your translation to propagate
// faster in the network you would apply more PoW, hence setting priority based on MWM.
type WeightQueue struct {
	q   pQueue
	c   chan bool
	mtx sync.Mutex
}

// NewWeightQueue creates a new weighted queue. Items in the queue are ordered by weight (heaviest first).
func NewWeightQueue(maxLen int) *WeightQueue {
	return &WeightQueue{
		q: make(pQueue, 0),
		c: make(chan bool, maxLen),
	}
}

// Push pushes a new weighted value to the queue. If the queue is full, this call blocks until the queue drops below
// its own max length.
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

// Pop pops an item from the queue. If no item is present the call blocks until a new item was pushed.
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
