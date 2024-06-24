package main

import (
	"container/heap"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"net"
	"github.com/fatih/color"
)

var replicaIPs []string
var ownIP string

var StateMachine map[string]int = make(map[string]int)

func main() {

	var seq int = 0
	var nullCnt int = 0

	//init SM

	//color output
	c := color.New(color.FgCyan)
	c.Add(color.Underline)

	// ログファイルを作成
	logFile, err := os.Create("logs/log.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer logFile.Close()

	//Choose the interval

	/*
		var interval int
		fmt.Println("Enter the interval: ")
		fmt.Scan(&interval)
	*/

	// サーバーに接続し、自身に割り当てられたポート番号を受け取る
	/*
		var Operation string
		fmt.Println("Operation: ")
		fmt.Scan(&Operation)/
	*/

	//Register to proxy
	ownIP = RegisterToProxy()

	// プロキシからの接続を待ち受ける
	// 他のレプリカのポート番号を取得
	replicaIPs = listenAndAcceptConnectionWithProxy()

	// 他のレプリカとの接続を確立
	go listenAndAccept()
	time.Sleep(250 * time.Millisecond)

	//ここで合意アルゴリズムを実行
	for {
		PQMutex.Lock()
		if PQ.Len() == 0 {
			PQMutex.Unlock()
			ConsensusTerminationMutex.Lock()
			if ConsensusTerminationMapList[seq] != nil {
				terminationValue := ConsensusTerminationMapList[seq][0]
				consensusValue := TerminationValue{isNull: terminationValue.Value == 0, CommandData: terminationValue.CommandData, phase: 0, seq: seq}
				notifyTermination(setConnectionWithOtherReplicas(replicaIPs),seq, consensusValue)
				color.Green("reached consensus: ", consensusValue, "\n")

				if !consensusValue.isNull{
					Dictionary[CommandTimestamp{Command: consensusValue.CommandData, Timestamp: consensusValue.CommandData.Timestamp}] = true
				}




				if !consensusValue.isNull && consensusValue.CommandData.Op == "" {
					 c := color.New(color.FgHiRed)
					 c.Println("This should not happen!")
				}
				if !consensusValue.isNull {
					parseWriteCommand(consensusValue.CommandData.Op, StateMachine)
				} else {
					nullCnt++
				}
				c.Println("SM in seq", seq, ":", StateMachine)

				IP2 := strings.Split(terminationValue.CommandData.ReplicaAddr, ":")[0]

				if ownIP == IP2 {
					// fmt.Println("Sending response to client")
					responseChannelMap[terminationValue.CommandData.ClientAddr] <- ResponseToClient{Value: 0, ClientAddr: terminationValue.CommandData.ClientAddr}
					// fmt.Println("Inserted response to slice")
				}

				seq++
				fmt.Println("Seq: ", seq)
			}

			ConsensusTerminationMutex.Unlock()
			continue
		}
		PQMutex.Unlock()
		commandPointer := heap.Pop(&PQ).(*CommandData)
		fmt.Println("Command: ", *commandPointer)
		if Dictionary[CommandTimestamp{Command: *commandPointer, Timestamp: commandPointer.Timestamp}] {
			//fmt.Println("Command already reached consensus: ", *commandPointer)
			//fmt.Println("Dictionary: ", Dictionary)
			delete(Dictionary, CommandTimestamp{Command: *commandPointer, Timestamp: commandPointer.Timestamp})
			continue
		}
		var stateStruct StateValueData
		// fmt.Println("cnt: ", seq)
		var terminationFlag int
		commandPointer.Seq = seq
		terminationFlag, stateStruct = exchangeStage(*commandPointer, seq)
		if terminationFlag == 1 {
			var consensusValue TerminationValue
			if stateStruct.Value == 0 {
				consensusValue = TerminationValue{isNull: true, CommandData: stateStruct.CommandData, phase: 0, seq: seq}
				notifyTermination(setConnectionWithOtherReplicas(replicaIPs),seq, consensusValue)
			}else{
				consensusValue = TerminationValue{isNull: false, CommandData: stateStruct.CommandData, phase: 0, seq: seq}
				notifyTermination(setConnectionWithOtherReplicas(replicaIPs),seq, consensusValue)
			}
			color.Green("reached consensus: ", consensusValue, "\n")

			if !consensusValue.isNull && consensusValue.CommandData.Op == "" {
				 c := color.New(color.FgHiRed)
				 c.Println("This should not happen!")
			}
			if consensusValue.CommandData != *commandPointer || consensusValue.isNull {
				PQMutex.Lock()
				// c := color.New(color.FgYellow)
				c.Println("Adding to dictionary: ", consensusValue.CommandData)
				PQ.Push(commandPointer)
				
				if !consensusValue.isNull {
					Dictionary[CommandTimestamp{Command: consensusValue.CommandData, Timestamp:consensusValue.CommandData.Timestamp}] = true
				}
				PQMutex.Unlock()
			}
			if !consensusValue.isNull {
				parseWriteCommand(consensusValue.CommandData.Op, StateMachine)
			} else {
				nullCnt++
			}

			//c.Println("SM in seq", seq, ":", StateMachine)

			// fmt.Println("IP: ", ownIP)
			IP2 := strings.Split(consensusValue.CommandData.ReplicaAddr, ":")[0]
			// fmt.Println("IP2", IP2)
			if ownIP == IP2 {
				// fmt.Println("Sending response to client")
				responseChannelMap[consensusValue.CommandData.ClientAddr] <- ResponseToClient{Value: 0, ClientAddr: consensusValue.CommandData.ClientAddr}
				// fmt.Println("Inserted response to slice")
			}

		

			seq++
			fmt.Println("Seq: ", seq)
			// fmt.Println("null cnt:", nullCnt)
			// fmt.Println("non-null percentage: ", (float64(seq-nullCnt)/float64(seq))*100)
			continue
		}

		consensusValue := weakMVC(stateStruct, seq)

		if !consensusValue.isNull && consensusValue.CommandData.Op == "" {
			 c := color.New(color.FgHiRed)
			 c.Println("This should not happen!")
		}
		if consensusValue.CommandData != *commandPointer || consensusValue.isNull {
			PQMutex.Lock()
			c := color.New(color.FgYellow)
			c.Println("Adding to dictionary: ", consensusValue.CommandData)
			PQ.Push(commandPointer)
			PQMutex.Unlock()
			if !consensusValue.isNull {
				Dictionary[CommandTimestamp{Command: consensusValue.CommandData, Timestamp: consensusValue.CommandData.Timestamp}] = true
			}
		}
		if !consensusValue.isNull {
			parseWriteCommand(consensusValue.CommandData.Op, StateMachine)
		} else {
			nullCnt++
		}
		// c.Println("SM in seq", seq, ":", StateMachine)

		IP2 := strings.Split(consensusValue.CommandData.ReplicaAddr, ":")[0]
		if ownIP == IP2 {
			 //fmt.Println("Sending response to client")
			 //fmt.Println("ClientAddr: ", consensusValue.CommandData.ClientAddr)
			responseChannelMap[consensusValue.CommandData.ClientAddr] <- ResponseToClient{Value: 0, ClientAddr: consensusValue.CommandData.ClientAddr}
			 //fmt.Println("Inserted response to slice")
		}

		//c.Println("SM in seq", seq, ":", StateMachine)
		seq++
		// fmt.Println("null cnt:", nullCnt)
		fmt.Println("Seq: ", seq)
		// fmt.Println("non-null percentage: ", (float64(seq-nullCnt)/float64(seq))*100)
	}
}

func weakMVC(stateStruct StateValueData, seq int) TerminationValue {

	var phase int = 0

	c := color.New(color.FgGreen)

	//Round 1
	//fmt.Println("State struct: ", stateStruct)
	var state StateValueData = StateValueData{Value: stateStruct.Value, Seq: seq, Phase: phase, CommandData: stateStruct.CommandData}
	terminationFlag, voteValue := roundOne(state, seq, phase)
	if terminationFlag == 1 {
		if voteValue.Value == 0 {
			terminationValue := TerminationValue{isNull: true, CommandData: voteValue.CommandData, phase: phase, seq: seq}
			notifyTermination(setConnectionWithOtherReplicas(replicaIPs),seq, terminationValue)
			c.Println("reached consensus: ", terminationValue)
			return terminationValue
		} else {
			terminationValue := TerminationValue{isNull: false, CommandData: voteValue.CommandData, phase: phase, seq: seq}
			notifyTermination(setConnectionWithOtherReplicas(replicaIPs),seq, terminationValue)
			c.Println("reached consensus: ", terminationValue)
			return terminationValue
		}
	}

	//Round 2
	//fmt.Println("voteValue: ", voteValue)
	var vote VoteValueData = VoteValueData{Value: voteValue.Value, Seq: seq, Phase: phase, CommandData: voteValue.CommandData}
	terminationFlag, returnStruct := roundTwo(vote, seq, phase)
	//fmt.Println("returnStruct: ", returnStruct)
	if terminationFlag == 1 {
		if returnStruct.ConsensusValue == 0 {
			terminationValue := TerminationValue{isNull: true, CommandData: returnStruct.CommandData, phase: phase, seq: seq}
			notifyTermination(setConnectionWithOtherReplicas(replicaIPs),seq, terminationValue)
			c.Println("eached consensus: ", terminationValue)
			return terminationValue
		} else {
			terminationValue := TerminationValue{isNull: false, CommandData: returnStruct.CommandData, phase: phase, seq: seq}
			notifyTermination(setConnectionWithOtherReplicas(replicaIPs),seq, terminationValue)
			c.Println("reached consensus: ", terminationValue)
			return terminationValue
		}
	}

	for {
		phase++

		
		state = StateValueData{Value: returnStruct.ConsensusValue, Seq: seq, Phase: phase, CommandData: returnStruct.CommandData}
		terminationFlag, voteValue = roundOne(state, seq, phase)
		if terminationFlag == 1 {
			if voteValue.Value == 0 {
				terminationValue := TerminationValue{isNull: true, CommandData: voteValue.CommandData, phase: phase, seq: seq}
				notifyTermination(setConnectionWithOtherReplicas(replicaIPs),seq, terminationValue)
				c.Println("reached consensus: ", terminationValue)
				return terminationValue
			} else {
				terminationValue := TerminationValue{isNull: false, CommandData: voteValue.CommandData, phase: phase, seq: seq}
				notifyTermination(setConnectionWithOtherReplicas(replicaIPs),seq, terminationValue)
				c.Println("reached consensus: ", terminationValue)
				return terminationValue
			}
		}
		//fmt.Println("voteValue: ", voteValue)
		var vote VoteValueData = VoteValueData{Value: voteValue.Value, Seq: seq, Phase: phase, CommandData: voteValue.CommandData}
		terminationFlag, returnStruct = roundTwo(vote, seq, phase)
		//fmt.Println("returnStruct: ", returnStruct)
		if terminationFlag == 1 {
			if returnStruct.ConsensusValue == 0 {
				terminationValue := TerminationValue{isNull: true, CommandData: returnStruct.CommandData, phase: phase, seq: seq}
				notifyTermination(setConnectionWithOtherReplicas(replicaIPs),seq, terminationValue)
				c.Println("reached consensus: ", terminationValue)
				return terminationValue
			} else {
				terminationValue := TerminationValue{isNull: false, CommandData: returnStruct.CommandData, phase: phase, seq: seq}
				notifyTermination(setConnectionWithOtherReplicas(replicaIPs),seq, terminationValue)
				c.Println("reached consensus: ", terminationValue)
				return terminationValue
			}
		}
		deleteData(seq, phase)
		c := color.New(color.FgHiRed)
		c.Println("No consensus reached in phase: ", phase)
	}
}


func notifyTermination(conns []net.Conn,  seq int, termination TerminationValue) {
	for _, conn := range conns {
		go func(conn net.Conn) {
			if termination.isNull{
				sendData(conn, ConsensusTermination{Seq: seq, Value: 0, CommandData: termination.CommandData})
			}else{
				sendData(conn, ConsensusTermination{Seq: seq, Value: 1, CommandData: termination.CommandData})
			}
		}(conn)
	}
}


