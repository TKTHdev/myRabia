package main

import (
    "fmt"
	"net"
	"sync"
	
)

var RoundOneMutex sync.Mutex
var RoundOneCntMutex sync.Mutex

func roundOne(state StateValueData, portNums []int, port int,seq int, phase int) int{
	conns := setConnectionWithOtherReplicas(portNums, port)
    wg := sync.WaitGroup{}
	roundOneSend(conns, state,phase, seq,&wg)
	var vote int =roundOneReceive(seq,len(portNums))
    return vote
}


func roundOneReceive(selfSeq int, nodeNum int) int {
    for{
		StateValueDataMutex.Lock()
        if(len(StateValueDataMapList[selfSeq])>=nodeNum/2+1){
            cnt := make(map[StateValueData]int)
            for _, command := range StateValueDataMapList[selfSeq] {
                RoundOneCntMutex.Lock()
                cnt[command]++
                RoundOneCntMutex.Unlock()
            }
            for _, c := range cnt {
                if c >= nodeNum/2+1 {
                    fmt.Println("State: 1")
					StateValueDataMutex.Unlock()
                    return 1
                }
            }
            fmt.Println("State: 0")
			StateValueDataMutex.Unlock()
            return 0
        }
		StateValueDataMutex.Unlock()
    }
}

func roundOneSend(conns []net.Conn, state StateValueData,phase int, seq int, wg *sync.WaitGroup) {
    for _, conn := range conns {
        wg.Add(1)
        go func(conn net.Conn) {
            defer wg.Done()
            sendData(conn, state)
        }(conn)
    }
}