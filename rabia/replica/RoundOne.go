package main

import (
	//"fmt"
	"net"
	"sync"
)

var RoundOneMutex sync.Mutex
var RoundOneCntMutex sync.Mutex

func roundOne(state StateValueData, seq int, phase int) (int, VoteValueData) {
	var returnCommand CommandData
	var voteValue int
	conns := setConnectionWithOtherReplicas(replicaIPs)
	wg := sync.WaitGroup{}
	roundOneSend(conns, state, &wg)
	terminationFlag, voteValue, returnCommand := roundOneReceive(seq, phase, len(portNums))
	returnVoteStruct := VoteValueData{Value: voteValue, Seq: seq, Phase: phase, CommandData: returnCommand}
	return terminationFlag, returnVoteStruct
}

func roundOneReceive(selfSeq int, phase int, nodeNum int) (int, int, CommandData) {
	for {
		var anyCommandReceived CommandData
		StateValueDataMutex.Lock()
		ConsensusTerminationMutex.Lock()
		if len(ConsensusTerminationMapList[selfSeq]) != 0 {
			value := ConsensusTerminationMapList[selfSeq][0].Value
			returnStruct := ConsensusTerminationMapList[selfSeq][0].CommandData
			ConsensusTerminationMutex.Unlock()
			StateValueDataMutex.Unlock()
			return 1, value, returnStruct
		}
		ConsensusTerminationMutex.Unlock()

		if (len(StateValueDataMapList[SeqPhase{Seq: selfSeq, Phase: phase}]) >= nodeNum/2+1) {
			for _, command := range StateValueDataMapList[SeqPhase{Seq: selfSeq, Phase: phase}] {
				if command.Value == 1 && command.CommandData.Op != "" {
					anyCommandReceived = command.CommandData
				}
			}
			//fmt.Println("any command received in state round: ", anyCommandReceived)
			cnt := make(map[StateValueData]int)
			for _, command := range StateValueDataMapList[SeqPhase{Seq: selfSeq, Phase: phase}] {
				RoundOneCntMutex.Lock()
				command.CommandData = anyCommandReceived
				cnt[command]++
				RoundOneCntMutex.Unlock()
			}
			for v, c := range cnt {
				if c >= nodeNum/2+1 {
					//fmt.Println("vote: ",v.Value)
					StateValueDataMutex.Unlock()
					return 0, v.Value, anyCommandReceived
				}
			}
			StateValueDataMutex.Unlock()
			//fmt.Println("vote: ?")
			return 0, -1, anyCommandReceived
		}
		StateValueDataMutex.Unlock()
	}
}

func roundOneSend(conns []net.Conn, state StateValueData, wg *sync.WaitGroup) {
	//fmt.Println("Sending state to replicas")
	for _, conn := range conns {
		wg.Add(1)
		go func(conn net.Conn) {
			defer wg.Done()
			sendData(conn, state)
		}(conn)
	}
}
