package main

import (
	//"fmt"
	"net"
	"sync"
)

var countMutex sync.Mutex

func exchangeStage(command CommandData, seq int) (int, StateValueData) {
	conns := setConnectionWithOtherReplicas(replicaIPs)
	var state int
	wg := sync.WaitGroup{}
	exchangeSend(conns, command, &wg)
	terminationFlag, state, CommandData := exchangeReceive(seq, len(replicaIPs))
	var returnStateValue StateValueData = StateValueData{Value: state, Seq: seq, CommandData: CommandData}
	//fmt.Println("State: ", state)
	return terminationFlag, returnStateValue
}

func exchangeReceive(selfSeq int, nodeNum int) (int, int, CommandData) {
	for {
		CommandDataMutex.Lock()
		ConsensusTerminationMutex.Lock()
		if len(ConsensusTerminationMapList[selfSeq]) != 0 {
			returnValue := ConsensusTerminationMapList[selfSeq][0].Value
			returnCommand := ConsensusTerminationMapList[selfSeq][0].CommandData
			CommandDataMutex.Unlock()
			ConsensusTerminationMutex.Unlock()
			return 1, returnValue, returnCommand
		}
		ConsensusTerminationMutex.Unlock()

		if len(CommandDataMapList[selfSeq]) >= nodeNum/2+1 {
			count := make(map[CommandData]int)
			for _, command := range CommandDataMapList[selfSeq] {
				countMutex.Lock()
				count[command]++
				countMutex.Unlock()
			}
			for v, c := range count {
				if c >= nodeNum/2+1 {
					CommandDataMutex.Unlock()
					return 0, 1, v
				}
			}
			CommandDataMutex.Unlock()
			return 0, 0, CommandData{}
		}
		CommandDataMutex.Unlock()
	}
}

func exchangeSend(conns []net.Conn, command CommandData, wg *sync.WaitGroup) {
	//fmt.Println("Sending Command: ", command)
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
