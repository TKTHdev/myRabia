package main

import (
	"container/heap"
	//"fmt"
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
		ConsensusTerminationMutex.Lock()
		if ConsensusTerminationMapList[seq] != nil {

			ConsensusTerminationMutex.Unlock()
			PQMutex.Unlock()


			value := ConsensusTerminationMapList[seq][0]
			resolveTermination(TerminationValue{isNull: value.Value==0, CommandData: value.CommandData, phase: 0, seq: seq}, CommandData{Op: "", Timestamp: time.Now(), ClientAddr: "", ReplicaAddr: ""})

			//c.Println("SM in seq", seq, ":", StateMachine)
			seq++
			continue

		}

		//fmt.Println("PQ: ", PQ)
		
		if len(PQ) == 0 {
			PQMutex.Unlock()
			ConsensusTerminationMutex.Unlock()
			continue
		}



		ConsensusTerminationMutex.Unlock()

		commandPointer := heap.Pop(&PQ).(*CommandData)


		if Dictionary[OpTimestamp{Op: commandPointer.Op, Timestamp: commandPointer.Timestamp}] {
			PQMutex.Unlock()
			continue
		}

		//fmt.Println("proposal: ", *commandPointer)
		PQMutex.Unlock()

		
		var stateStruct StateValueData
		// fmt.Println("cnt: ", seq)
		var terminationFlag int
		commandPointer.Seq = seq
		terminationFlag, stateStruct = exchangeStage(*commandPointer, seq)
		if terminationFlag == 1 {

			value := TerminationValue{isNull: stateStruct.Value == 0, CommandData: stateStruct.CommandData, phase: 0, seq: seq}
			//color.Green("reached consensus: ", value, "\n")
			resolveTermination(value, *commandPointer)
			//c.Println("SM in seq", seq, ":", StateMachine)
			seq++
			// fmt.Println("null cnt:", nullCnt)
			// fmt.Println("non-null percentage: ", (float64(seq-nullCnt)/float64(seq))*100)
			continue
		}

		consensusValue := weakMVC(stateStruct, seq)
		resolveTermination(consensusValue, *commandPointer)


		//c.Println("SM in seq", seq, ":", StateMachine)
		seq++
		// fmt.Println("null cnt:", nullCnt)
		// fmt.Println("non-null percentage: ", (float64(seq-nullCnt)/float64(seq))*100)
	}
}

func weakMVC(stateStruct StateValueData, seq int) TerminationValue {

	var phase int = 0

	//c := color.New(color.FgGreen)

	//Round 1
	//fmt.Println("State struct: ", stateStruct)
	var state StateValueData = StateValueData{Value: stateStruct.Value, Seq: seq, Phase: phase, CommandData: stateStruct.CommandData}
	terminationFlag, voteValue := roundOne(state, seq, phase)
	if terminationFlag == 1 {
		if voteValue.Value == 0 {
			terminationValue := TerminationValue{isNull: true, CommandData: voteValue.CommandData, phase: phase, seq: seq}
			notifyTermination(setConnectionWithOtherReplicas(replicaIPs),seq, terminationValue)
		//	c.Println("reached consensus: ", terminationValue)
			return terminationValue
		} else {
			terminationValue := TerminationValue{isNull: false, CommandData: voteValue.CommandData, phase: phase, seq: seq}
			notifyTermination(setConnectionWithOtherReplicas(replicaIPs),seq, terminationValue)
			//c.Println("reached consensus: ", terminationValue)
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
			//c.Println("eached consensus: ", terminationValue)
			return terminationValue
		} else {
			terminationValue := TerminationValue{isNull: false, CommandData: returnStruct.CommandData, phase: phase, seq: seq}
			notifyTermination(setConnectionWithOtherReplicas(replicaIPs),seq, terminationValue)
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
				notifyTermination(setConnectionWithOtherReplicas(replicaIPs),seq, terminationValue)
				//c.Println("reached consensus: ", terminationValue)
				return terminationValue
			} else {
				terminationValue := TerminationValue{isNull: false, CommandData: voteValue.CommandData, phase: phase, seq: seq}
				notifyTermination(setConnectionWithOtherReplicas(replicaIPs),seq, terminationValue)
				//c.Println("reached consensus: ", terminationValue)
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
				//c.Println("reached consensus: ", terminationValue)
				return terminationValue
			} else {
				terminationValue := TerminationValue{isNull: false, CommandData: returnStruct.CommandData, phase: phase, seq: seq}
				notifyTermination(setConnectionWithOtherReplicas(replicaIPs),seq, terminationValue)
				//c.Println("reached consensus: ", terminationValue)
				return terminationValue
			}
		}
		deleteData(seq, phase)
		//c := color.New(color.FgHiRed)
		//c.Println("No consensus reached in phase: ", phase)
	}
}


func notifyTermination(conns []net.Conn,  seq int, termination TerminationValue) {
	for _, conn := range conns {
		go func(conn net.Conn) {
			if termination.isNull{
				sendData(conn, ConsensusTermination{Seq: seq, Value: 0, CommandData: termination.CommandData})
				//fmt.Println("sending termination to: ", conn.RemoteAddr().String())
			}else{
				sendData(conn, ConsensusTermination{Seq: seq, Value: 1, CommandData: termination.CommandData})
				//fmt.Println("sending termination to: ", conn.RemoteAddr().String())
			}
		}(conn)
	}
}


func resolveTermination(termination TerminationValue, ownProposal CommandData){
	if !termination.isNull && termination.CommandData.Op == "" {
		//c := color.New(color.FgHiRed)
		//c.Println("This should not happen!")
	}

	if termination.isNull && ownProposal.Op != ""{
		PQMutex.Lock()
		heap.Push(&PQ, &ownProposal)
		PQMutex.Unlock()
		return 
	}

	if termination.CommandData != ownProposal{
		PQMutex.Lock()
		if ownProposal.Op != "" {
			heap.Push(&PQ, &ownProposal)
		}
		PQMutex.Unlock()
		Dictionary[OpTimestamp{Op: termination.CommandData.Op, Timestamp: termination.CommandData.Timestamp}] = true
	}
	parseWriteCommand(termination.CommandData.Op, StateMachine)

	IP := strings.Split(termination.CommandData.ReplicaAddr, ":")[0]
	if ownIP == IP  {
		//c:= color.New(color.FgHiMagenta)
		//c.Println("sending response to client: ", termination.CommandData.ClientAddr)

		responseChannelMap[termination.CommandData.ClientAddr] <- ResponseToClient{Value: 0, ClientAddr: termination.CommandData.ClientAddr}
	}
}