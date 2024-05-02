package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
	"github.com/fatih/color"
	
)

var listener net.Listener
var portNums []int

func main() {

	// ログファイルを作成
	var logName string 
	fmt.Println("Enter log file name: ")
	fmt.Scan(&logName)
	logName = "log" + logName + ".txt"
	logFile, err := os.Create("logs/"+logName)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	
	// サーバーに接続し、自身に割り当てられたポート番号を受け取る
	/*
	var Operation string
	fmt.Println("Operation: ")
	fmt.Scan(&Operation)
	*/
	var port string = sayHelloAndReceivePortNum()
	portInt, _ := strconv.Atoi(port)

	// プロキシからの接続を待ち受ける
	// 他のレプリカのポート番号を取得
	portNums = listenAndAcceptConnectionWithProxy(listener, port)

	// 他のレプリカとの接続を確立
	go listenAndAccept(port)	
	time.Sleep(250 * time.Millisecond)


	//ここで合意アルゴリズムを実行
	var seq int = 0
	for {

		PQMutex.Lock()
		commandPointer :=PQ.Head()
		if commandPointer == nil{
			PQMutex.Unlock()
			continue
		}
		PQ.Pop()
		PQMutex.Unlock()
		var stateStruct StateValueData
		fmt.Println("cnt: ", seq)
		var terminationFlag int
		commandPointer.Seq = seq
		terminationFlag, stateStruct = exchangeStage(*commandPointer, portNums, seq)
		if terminationFlag == 1{
			logger.Println("consensusValue: ", stateStruct.Value, "command: ", stateStruct.CommandData.Op)
			terminationValue := TerminationValue{isNull: false, CommandData: stateStruct.CommandData, phase: 0, seq: seq}
			color.Green("reached consensus: ", terminationValue,"\n")
			seq++
			continue
		}
		

	    consensusValue:=weakMVC(stateStruct, portInt, portNums, listener, seq)
		logger.Println("consensusValue: ", consensusValue)
		seq++
		
		//delete data to save memory
		time.Sleep(250 * time.Millisecond)

	}

}

func weakMVC(stateStruct StateValueData , selfPort int, portNums []int, ln net.Listener, seq int) (TerminationValue){

	var phase int = 0

	//Round 1
	fmt.Println("State struct: ", stateStruct)
	var state StateValueData = StateValueData{Value: stateStruct.Value, Seq: seq, Phase: phase, CommandData: stateStruct.CommandData}
	terminationFlag, voteValue :=roundOne(state, portNums,  seq, phase)
	if terminationFlag == 1{
		
		if voteValue.Value == 0{
			terminationValue := TerminationValue{isNull: true, CommandData: voteValue.CommandData, phase: phase, seq: seq}
			color.Green("reached consensus: ", terminationValue,"\n")
			return terminationValue
		}else{
			terminationValue := TerminationValue{isNull: false, CommandData: voteValue.CommandData, phase: phase, seq: seq}
			color.Green("reached consensus: ", terminationValue,"\n")
			return terminationValue
		}
	}	
	//Round 2
	fmt.Println("voteValue: ", voteValue)
	var vote VoteValueData = VoteValueData{Value: voteValue.Value, Seq: seq, Phase: phase, CommandData: voteValue.CommandData}
	terminationFlag, returnStruct :=roundTwo(vote, portNums, selfPort, seq,phase)
	fmt.Println("returnStruct: ", returnStruct)
	if(terminationFlag == 1){
		if returnStruct.ConsensusValue == 0{
			terminationValue := TerminationValue{isNull: true, CommandData: returnStruct.CommandData, phase: phase, seq: seq}
			color.Green("reached consensus: ", terminationValue,"\n")
			return terminationValue
		}else{
			terminationValue := TerminationValue{isNull: false, CommandData: returnStruct.CommandData, phase: phase, seq: seq}
			color.Green("reached consensus: ", terminationValue,"\n")
			return terminationValue
		}
	}	

	for{
		phase++

		state = StateValueData{Value: returnStruct.ConsensusValue, Seq: seq, Phase: phase, CommandData: returnStruct.CommandData}
		terminationFlag, voteValue =roundOne(state, portNums,  seq, phase)
		if terminationFlag == 1{
			if voteValue.Value == 0{
				terminationValue := TerminationValue{isNull: true, CommandData: voteValue.CommandData, phase: phase, seq: seq}
				color.Green("reached consensus: ", terminationValue,"\n")
				return terminationValue
			}else{
				terminationValue := TerminationValue{isNull: false, CommandData: voteValue.CommandData, phase: phase, seq: seq}
				color.Green("reached consensus: ", terminationValue,"\n")
				return terminationValue
			}
		}
		fmt.Println("voteValue: ", voteValue)
		var vote  VoteValueData = VoteValueData{Value: voteValue.Value, Seq: seq, Phase: phase, CommandData: voteValue.CommandData}
		terminationFlag,returnStruct =roundTwo(vote, portNums, selfPort, seq,phase)
		fmt.Println("returnStruct: ", returnStruct)
		if terminationFlag == 1{
			if returnStruct.ConsensusValue == 0{
				terminationValue := TerminationValue{isNull: true, CommandData: returnStruct.CommandData, phase: phase, seq: seq}
				color.Green("reached consensus: ", terminationValue,"\n")
			}else{
				terminationValue := TerminationValue{isNull: false, CommandData: returnStruct.CommandData, phase: phase, seq: seq}
				color.Green("reached consensus: ", terminationValue,"\n")
				return terminationValue
			}
		}
		deleteData(seq, phase)
	}
}
