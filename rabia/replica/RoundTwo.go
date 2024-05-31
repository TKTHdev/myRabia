package main

import (
	"net"
	//"fmt"
	"math/rand"
	"sync"
)

var RoundTwoMutex sync.Mutex
var RoundTwoCntMutex sync.Mutex

func roundTwo(vote VoteValueData, seq int, phase int) (int, RoundTwoReturnStruct) {
	conns := setConnectionWithOtherReplicas(replicaIPs)
	var returnStruct RoundTwoReturnStruct
	var terminationValue int
	wg := sync.WaitGroup{}
	roundTwoSend(conns, vote, &wg)
	terminationFlag, terminationValue, CommandData := roundTwoReceive(seq, len(replicaIPs), phase)
	if CommandData.Op != "" {
		CommandData = vote.CommandData
	}
	returnStruct = RoundTwoReturnStruct{ConsensusValue: terminationValue, CommandData: CommandData}
	return terminationFlag, returnStruct
}

func roundTwoReceive(selfSeq int, nodeNum int, phase int ) (int, int, CommandData) {
	for {
		VoteValueDataMutex.Lock()
		var anyCommandReceived CommandData
		ConsensusTerminationMutex.Lock()
		if len(ConsensusTerminationMapList[selfSeq]) != 0 {
			//fmt.Println(len(ConsensusTerminationMapList[selfSeq]))
			value := ConsensusTerminationMapList[selfSeq][0].Value
			ConsensusTerminationMutex.Unlock()
			VoteValueDataMutex.Unlock()
			return 1, value, ConsensusTerminationMapList[selfSeq][0].CommandData
		}
		ConsensusTerminationMutex.Unlock()

		if (len(VoteValueDataMapList[SeqPhase{Seq: selfSeq, Phase: phase}]) >= nodeNum/2+1) {
			for _, command := range VoteValueDataMapList[SeqPhase{Seq: selfSeq, Phase: phase}] {
				if command.CommandData.Op != "" {
					anyCommandReceived = command.CommandData
				}
			}
			//fmt.Println("VoteValueDataMapList: ", VoteValueDataMapList[SeqPhase{Seq: selfSeq, Phase: phase}])
			//fmt.Println("any command received in vote round: ", anyCommandReceived)
			cnt := make(map[VoteValueData]int)
			for _, command := range VoteValueDataMapList[SeqPhase{Seq: selfSeq, Phase: phase}] {
				RoundTwoCntMutex.Lock()
				command.CommandData = anyCommandReceived
				cnt[command]++

				RoundTwoCntMutex.Unlock()
			}

			for v, c := range cnt {
				if c >= nodeNum/2+1 && v.Value != -1 {
					VoteValueDataMutex.Unlock()
					notifyTermination(setConnectionWithOtherReplicas(replicaIPs), &sync.WaitGroup{}, selfSeq, v)
					return 1, v.Value, anyCommandReceived
				}
			}
			for v, c := range cnt {
				if c >= 1 && v.Value != -1 {
					//fmt.Println("found at least one vote for non-? value: ",v.Value)
					VoteValueDataMutex.Unlock()
					return 0, v.Value, anyCommandReceived
				}
			}
			stateCoinFlip := CommonCoinFlip(selfSeq, phase)
			//fmt.Println("coin flip: ",stateCoinFlip)
			VoteValueDataMutex.Unlock()
			return 0, stateCoinFlip, anyCommandReceived
		}
		VoteValueDataMutex.Unlock()
	}
}

func roundTwoSend(conns []net.Conn, vote VoteValueData, wg *sync.WaitGroup) {
	for _, conn := range conns {
		wg.Add(1)
		go func(conn net.Conn) {
			//fmt.Println("sending:  ", vote)
			defer wg.Done()
			sendData(conn, vote)
		}(conn)
	}
}

func notifyTermination(conns []net.Conn, wg *sync.WaitGroup, seq int, termination VoteValueData) {
	for _, conn := range conns {
		wg.Add(1)
		go func(conn net.Conn) {

			defer wg.Done()
			sendData(conn, ConsensusTermination{Seq: seq, Value: termination.Value, CommandData: termination.CommandData})
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
