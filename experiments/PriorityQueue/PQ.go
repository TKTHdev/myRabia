package main

import (
    "container/heap"
    "math/rand"
)

type CommandData struct {
    Op        string
    Timestamp int
    Seq       int
}

type PriorityQueue []*CommandData

func (pq PriorityQueue) Len() int {
    return len(pq)
}

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
    item := old[n-1]
    *pq = old[0 : n-1]
    return item
}

func generateRandomCommand() string {
    commands := []string{"read", "write", "delete"}
    return commands[rand.Intn(len(commands))]
}

func CreateRandomCommand() *CommandData {
    return &CommandData{Op: generateRandomCommand(), Timestamp: 0, Seq: 0}
}

func main() {
    var pq PriorityQueue
    pq = make(PriorityQueue, 0)

    for i := 10; i < 50; i++ {
        command := CreateRandomCommand()
        command.Timestamp = i
        command.Seq = i
        heap.Push(&pq, command)
    }

    for i := 0; i < 10; i++ {
        command := CreateRandomCommand()
        command.Timestamp = i
        command.Seq = i
        heap.Push(&pq, command)
    }

    for i := 50; i < 100; i++ {
        command := CreateRandomCommand()
        command.Timestamp = i
        command.Seq = i
        heap.Push(&pq, command)
    }

    for pq.Len() > 0 {
        command := heap.Pop(&pq).(*CommandData)
        println(command.Timestamp)
    }
}