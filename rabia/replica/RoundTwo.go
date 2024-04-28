package main

import (
    "fmt"
	"net"
	"sync"
	
)

var RoundTwoMutex sync.Mutex
var RoundTwoCntMutex sync.Mutex

func roundTwo(vote VoteValueData, portNums []int, port int,seq int, phase int) int{
	conns := setConnectionWithOtherReplicas(portNums, port)
    wg := sync.WaitGroup{}
	roundTwoSend(conns, vote,phase, seq,&wg)
	var consensus int =roundTwoReceive(seq,len(portNums))
    return consensus
}


func roundTwoReceive(selfSeq int, nodeNum int) int {
    for{
		VoteValueDataMutex.Lock()
        if(len(VoteValueDataMapList[selfSeq])>=nodeNum/2+1){
            cnt := make(map[VoteValueData]int)
            for _, command := range VoteValueDataMapList[selfSeq] {
                RoundTwoCntMutex.Lock()
                cnt[command]++
                RoundTwoCntMutex.Unlock()
            }
            for _, c := range cnt {
                if c >= nodeNum/2+1 {
                    fmt.Println("State: 1")
					VoteValueDataMutex.Unlock()
                    return 1
                }
            }
            fmt.Println("State: 0")
			VoteValueDataMutex.Unlock()
            return 0
        }
		VoteValueDataMutex.Unlock()
    }
}

func roundTwoSend(conns []net.Conn, vote VoteValueData,phase int, seq int, wg *sync.WaitGroup) {
    for _, conn := range conns {
        wg.Add(1)
        go func(conn net.Conn) {
            defer wg.Done()
            sendData(conn, vote)
        }(conn)
    }
}