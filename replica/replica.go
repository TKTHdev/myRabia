package main

import (
	"container/heap"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

var replicaIPs []string
var ownIP string

var StateMachine map[string]int = make(map[string]int)


func main() {

	var nullCnt int = 0
	var seq int = 0

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
	logger := log.New(logFile, "", log.LstdFlags)


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
			continue
		}
		commandPointer := heap.Pop(&PQ).(*CommandData)
		logger.Println("Command: ", *commandPointer, "Dict: ", Dictionary)
		if Dictionary[CommandTimestamp{Command: commandPointer.Op, Timestamp: commandPointer.Timestamp}] {
			PQMutex.Unlock()
			fmt.Println("Command already reached consensus: ", *commandPointer)
			delete(Dictionary, CommandTimestamp{Command: commandPointer.Op, Timestamp: commandPointer.Timestamp})
			continue
		}
		logger.Println("Proposing command: ", *commandPointer)
		PQMutex.Unlock()
		var stateStruct StateValueData
		fmt.Println("cnt: ", seq)
		var terminationFlag int
		commandPointer.Seq = seq
		terminationFlag, stateStruct = exchangeStage(*commandPointer, seq)
		if terminationFlag == 1 {
			consensusValue := TerminationValue{isNull: false, CommandData: stateStruct.CommandData, phase: 0, seq: seq}
			//color.Green("reached consensus: ", consensusValue, "\n")

			if !consensusValue.isNull && consensusValue.CommandData.Op == "" {
				c := color.New(color.FgHiRed)
				c.Println("This should not happen!")
			}
			if consensusValue.CommandData != *commandPointer || consensusValue.isNull {
				PQMutex.Lock()
				c := color.New(color.FgYellow)
				c.Println("Adding to dictionary: ", consensusValue.CommandData)
				PQ.Push(commandPointer)
				Dictionary[CommandTimestamp{Command: consensusValue.CommandData.Op, Timestamp: consensusValue.CommandData.Timestamp}] = true
				PQMutex.Unlock()
			}
			if !consensusValue.isNull {
				parseWriteCommand(consensusValue.CommandData.Op, StateMachine)
			}

			if consensusValue.isNull {
				nullCnt++
			} 
			c.Println("SM in seq", seq, ":", StateMachine)

			IP2 := strings.Split(consensusValue.CommandData.ReplicaAddr, ":")[0]
			if ownIP == IP2 {
				terminationChannelMutex.Lock()
				responseSlice = append(responseSlice, ResponseToClient{Value: 0, ClientAddr: consensusValue.CommandData.ClientAddr})
				terminationChannelMutex.Unlock()
			}

			//Print the size of PQ
			seq++
			report(nullCnt, seq)
		}

		consensusValue := weakMVC(stateStruct, seq)

		logger.Println("consensusValue: ", consensusValue)
		if !consensusValue.isNull && consensusValue.CommandData.Op == "" {
			c := color.New(color.FgHiRed)
			c.Println("This should not happen!")
		}
		if consensusValue.CommandData != *commandPointer || consensusValue.isNull {
			PQMutex.Lock()
			c := color.New(color.FgYellow)
			c.Println("Adding to dictionary: ", consensusValue.CommandData)
			PQ.Push(commandPointer)
			Dictionary[CommandTimestamp{Command: consensusValue.CommandData.Op, Timestamp: consensusValue.CommandData.Timestamp}] = true
			PQMutex.Unlock()
		}
		if !consensusValue.isNull {
			parseWriteCommand(consensusValue.CommandData.Op, StateMachine)
		}
		c.Println("SM in seq", seq, ":", StateMachine)

		IP2 := strings.Split(consensusValue.CommandData.ReplicaAddr, ":")[0]
		if ownIP == IP2 {
			terminationChannelMutex.Lock()
			responseSlice = append(responseSlice, ResponseToClient{Value: 0, ClientAddr: consensusValue.CommandData.ClientAddr})
			terminationChannelMutex.Unlock()
		}

		//Print the size of PQ
		seq++
		report(nullCnt, seq)

	}

}

func weakMVC(stateStruct StateValueData, seq int) TerminationValue {

	var phase int = 0

	//c := color.New(color.FgGreen)

	//Round 1
	var state StateValueData = StateValueData{Value: stateStruct.Value, Seq: seq, Phase: phase, CommandData: stateStruct.CommandData}
	terminationFlag, voteValue := roundOne(state, seq, phase)
	if terminationFlag == 1 {

		if voteValue.Value == 0 {
			terminationValue := TerminationValue{isNull: true, CommandData: voteValue.CommandData, phase: phase, seq: seq}
			//c.Println("reached consensus: ", terminationValue)
			return terminationValue
		} else {
			terminationValue := TerminationValue{isNull: false, CommandData: voteValue.CommandData, phase: phase, seq: seq}
			//c.Println("reached consensus: ", terminationValue)
			return terminationValue
		}
	}

	//Round 2
	var vote VoteValueData = VoteValueData{Value: voteValue.Value, Seq: seq, Phase: phase, CommandData: voteValue.CommandData}
	terminationFlag, returnStruct := roundTwo(vote, seq, phase)
	//fmt.Println("returnStruct: ", returnStruct)
	if terminationFlag == 1 {
		if returnStruct.ConsensusValue == 0 {
			terminationValue := TerminationValue{isNull: true, CommandData: returnStruct.CommandData, phase: phase, seq: seq}
			//c.Println("reached consensus: ", terminationValue)
			return terminationValue
		} else {
			terminationValue := TerminationValue{isNull: false, CommandData: returnStruct.CommandData, phase: phase, seq: seq}
			//c.Println("reached consensus: ", terminationValue)
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
				//c.Println("reached consensus: ", terminationValue)
				return terminationValue
			} else {
				terminationValue := TerminationValue{isNull: false, CommandData: voteValue.CommandData, phase: phase, seq: seq}
				//c.Println("reached consensus: ", terminationValue)
				return terminationValue
			}
		}
		var vote VoteValueData = VoteValueData{Value: voteValue.Value, Seq: seq, Phase: phase, CommandData: voteValue.CommandData}
		terminationFlag, returnStruct = roundTwo(vote, seq, phase)
		//fmt.Println("returnStruct: ", returnStruct)
		if terminationFlag == 1 {
			if returnStruct.ConsensusValue == 0 {
				terminationValue := TerminationValue{isNull: true, CommandData: returnStruct.CommandData, phase: phase, seq: seq}
				//c.Println("reached consensus: ", terminationValue)
				return terminationValue
			} else {
				terminationValue := TerminationValue{isNull: false, CommandData: returnStruct.CommandData, phase: phase, seq: seq}
				//c.Println("reached consensus: ", terminationValue)
				return terminationValue
			}
		}
		deleteData(seq, phase)
	}
}

func report(nullCnt int, seq int){
	fmt.Println("nullCount: ", nullCnt, "Percentage of non-null consensus: ", (1-float64(nullCnt)/float64(seq))*100, "%")
}
