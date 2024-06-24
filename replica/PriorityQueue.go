package main


import(
	"container/heap"
)

type PriorityQueue []*CommandData

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	
	return pq[i].Timestamp.Before(pq[j].Timestamp)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*CommandData)
	*pq = append(*pq, item)
	heap.Fix(pq, pq.Len()-1)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
