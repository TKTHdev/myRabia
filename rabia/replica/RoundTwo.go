package main

import (
	"net"
    "fmt"
	"sync"
    "math/rand"
	
)

var RoundTwoMutex sync.Mutex
var RoundTwoCntMutex sync.Mutex

func roundTwo(vote VoteValueData, portNums []int, port int,seq int, phase int) (int,int){
	conns := setConnectionWithOtherReplicas(portNums)
    wg := sync.WaitGroup{}
	roundTwoSend(conns, vote,phase, seq,&wg)
	var consensus,consensusValue  =roundTwoReceive(seq,len(portNums),phase)
    return consensus,consensusValue
}


func roundTwoReceive(selfSeq int, nodeNum int,phase int) (int,int) {
    for{
		VoteValueDataMutex.Lock()
        if(len(VoteValueDataMapList[selfSeq])>=nodeNum/2+1){
            cnt := make(map[VoteValueData]int)
            for _, command := range VoteValueDataMapList[selfSeq] {
                RoundTwoCntMutex.Lock()
                cnt[command]++
                RoundTwoCntMutex.Unlock()
            }

           
            for v, c := range cnt {
                if c >= nodeNum/2+1 && v.Value != -1{
					VoteValueDataMutex.Unlock()
                    return -1,v.Value
                }
            }
            for v, c := range cnt {
                if c>=1 && v.Value != -1{
                    fmt.Println("found at least one vote for non-? value: ",v.Value)
                    VoteValueDataMutex.Unlock()
                    return v.Value, -1
                }
            }
            stateCoinFlip :=CommonCoinFlip(selfSeq, phase)
            fmt.Println("coin flip: ",stateCoinFlip)
			VoteValueDataMutex.Unlock()
            return stateCoinFlip,-1 
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

func CommonCoinFlip(seq, phase int) int {
    seed := int64(seq*1000 + phase)
    rand.Seed(seed)

    // 乱数が0.5以上の場合は1を、そうでない場合は0を返す
    if rand.Float64() >= 0.5 {
        return 1
    } else {
        return 0
    }
}
