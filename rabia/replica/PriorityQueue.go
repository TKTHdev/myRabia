package main




type PriorityQueue []*CommandData





func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
    return pq[i].Timestamp < pq[j].Timestamp
}

func (pq PriorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
    item := x.(*CommandData)
    *pq = append(*pq, item)
}


func (pq *PriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    if n == 0 {
        return nil
    }
    item := old[0]
    *pq = old[1:n]
    return item
}

func (pq *PriorityQueue) Head() *CommandData {
    if pq.Len() == 0 {
        return nil
    }
    return (*pq)[0]
}