package main

import (
    //"fmt"
	"net"
	"sync"
	
)

var countMutex sync.Mutex

func exchangeStage(command CommandData, portNums []int, seq int) int{
	conns := setConnectionWithOtherReplicas(portNums)
    var state int
    wg := sync.WaitGroup{}
	exchangeSend(conns, command, &wg)
	state = exchangeReceive(seq,len(portNums))
    //fmt.Println("State: ", state)
    return state 
}


func exchangeReceive(selfSeq int, nodeNum int) int {
    for{
        CommandDataMutex.Lock()
        if(len(CommandDataMapList[selfSeq])>=nodeNum/2+1){
            count := make(map[CommandData]int)
            for _, command := range CommandDataMapList[selfSeq] {
                countMutex.Lock()
                count[command]++
                countMutex.Unlock()
            }
            for _, c := range count {
                if c >= nodeNum/2+1 {
                    CommandDataMutex.Unlock()
                    return 1
                }
            }
            CommandDataMutex.Unlock()
            return 0
        }
        CommandDataMutex.Unlock()
    }
}

func exchangeSend(conns []net.Conn, command CommandData, wg *sync.WaitGroup) {
    for _, conn := range conns {
        wg.Add(1)
        go func(conn net.Conn) {
            defer wg.Done()
            //fmt.Printf("Sending command to %v\n", conn.RemoteAddr())
            //fmt.Println("Sending Command: ", command)
            sendData(conn, command)
        }(conn)
    }
}