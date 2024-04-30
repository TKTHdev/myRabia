package main

import (
    "fmt"
	"net"
	"sync"
	
)

var RoundOneMutex sync.Mutex
var RoundOneCntMutex sync.Mutex

func roundOne(state StateValueData, portNums []int, seq int, phase int) (int,int){
	conns := setConnectionWithOtherReplicas(portNums)
    wg := sync.WaitGroup{}
	roundOneSend(conns, state,&wg)
    terminationFlag,vote :=roundOneReceive(seq, phase, len(portNums))
    return terminationFlag, vote
}


func roundOneReceive(selfSeq int, phase int,  nodeNum int) (int,int) {
    for{
		StateValueDataMutex.Lock()
        ConsensusTerminationMutex.Lock()
        if len(ConsensusTerminationMapList[selfSeq]) !=0 {
            value := ConsensusTerminationMapList[selfSeq][0].Value
            ConsensusTerminationMutex.Unlock()
            StateValueDataMutex.Unlock()
            return 1,value
        }
        ConsensusTerminationMutex.Unlock()
        
        if(len(StateValueDataMapList[SeqPhase{Seq: selfSeq, Phase: phase}])>=nodeNum/2+1){
            cnt := make(map[StateValueData]int)
            for _, command := range StateValueDataMapList[SeqPhase{ Seq: selfSeq, Phase: phase }] {
                RoundOneCntMutex.Lock()
                cnt[command]++
                RoundOneCntMutex.Unlock()
            }
            for v, c := range cnt {
                if c >= nodeNum/2+1 {
                    //fmt.Println("vote: ",v.Value)
					StateValueDataMutex.Unlock()
                    return 0,v.Value
                }
            }
			StateValueDataMutex.Unlock()
            //fmt.Println("vote: ?")
            return 0,-1
        }
		StateValueDataMutex.Unlock()
    }
}

func roundOneSend(conns []net.Conn, state StateValueData, wg *sync.WaitGroup) {
    fmt.Println("Sending state to replicas")
    for _, conn := range conns {
        wg.Add(1)
        go func(conn net.Conn) {
            defer wg.Done()
            sendData(conn, state)
        }(conn)
    }
}