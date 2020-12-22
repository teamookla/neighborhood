package neighborhood

import "container/heap"

type item struct {
	point    Point
	node     *kdTreeNode
	distance float64
	rank     float64
}

// priorityQueue implements heap.Interface and holds searchPoints
type priorityQueue []*item

// newPriorityQueue creates a rank queue a given initial capacity
func newPriorityQueue(capacity int) priorityQueue {
	queue := make(priorityQueue, 0, capacity)
	heap.Init(&queue)
	return queue
}

// PushPoint creates a new Point item and pushes it into the queue
func (pq *priorityQueue) PushPoint(point Point, dist float64) {
	// see if point implements optional Ranker interface
	rank := 0.0
	if ranked, ok := point.(Ranker); ok {
		rank = ranked.GetRank()
	}
	heap.Push(pq, &item{
		point:    point,
		distance: dist,
		rank:     rank,
	})
}

// PushPoint creates a new searchPoint and pushes it into the queue
func (pq *priorityQueue) PushNode(node *kdTreeNode) {
	// see if point implements optional Ranker interface
	heap.Push(pq, &item{
		node:     node,
		distance: node.Dist,
		rank:     -1.0,
	})
}

func (pq *priorityQueue) PopItem() *item {
	if i := heap.Pop(pq); i != nil {
		return i.(*item)
	}
	return nil
}

func (pq *priorityQueue) Peek() (itm *item) {
	items := *pq
	if len(items) > 0 {
		itm = items[0]
	}
	return
}

//
// heap.Interface implementation
//

func (pq priorityQueue) Less(i, j int) bool {
	// check for equal distances
	if pq[i].distance == pq[j].distance {
		// Pop the highest rank (tie breaker)
		return pq[i].rank > pq[j].rank
	}
	// Pop the lowest distance
	return pq[i].distance < pq[j].distance
}

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Swap(i, j int) {
	if len(pq) > i && len(pq) > j {
		pq[i], pq[j] = pq[j], pq[i]
	}
}

func (pq *priorityQueue) Push(x interface{}) {
	itm := x.(*item)
	*pq = append(*pq, itm)
}

func (pq *priorityQueue) Pop() interface{} {
	if pq == nil || len(*pq) < 1 {
		return nil
	}
	old := *pq
	n := len(old)
	i := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return i
}
