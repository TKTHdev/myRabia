package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
	"github.com/fatih/color"
	"container/heap"
)

var listener net.Listener
var portNums []int

func main() {

	//init SM
	StateMachine := make(map[string]int)

	//color output
	c := color.New(color.FgCyan)
	c.Add(color.Underline)

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
		if PQ.Len() == 0{
			PQMutex.Unlock()
			continue
		}
		commandPointer := heap.Pop(&PQ).(*CommandData)
		if Dictionary[*commandPointer]{
			PQMutex.Unlock()
			Dictionary[*commandPointer] = false
			delete(Dictionary, *commandPointer)
			continue
		}
		PQMutex.Unlock()
		var stateStruct StateValueData
		fmt.Println("cnt: ", seq)
		var terminationFlag int
		commandPointer.Seq = seq
		terminationFlag, stateStruct = exchangeStage(*commandPointer, portNums, seq)
		if terminationFlag == 1{
			terminationValue := TerminationValue{isNull: false, CommandData: stateStruct.CommandData, phase: 0, seq: seq}
			logger.Println("consensusValue: ", terminationValue)
			color.Green("reached consensus: ", terminationValue,"\n")
			seq++
			continue
		}
		
	    consensusValue:=weakMVC(stateStruct, portInt, portNums,seq)
		logger.Println("consensusValue: ", consensusValue)
		if !consensusValue.isNull && consensusValue.CommandData.Op == "" {
			c := color.New(color.FgHiRed)
			c.Println("This should not happen!")
		}
		if consensusValue.CommandData != *commandPointer || consensusValue.isNull {
			PQMutex.Lock()
			fmt.Println("Not Matching:")
			heap.Push(&PQ, commandPointer)
			Dictionary[consensusValue.CommandData] = true
			PQMutex.Unlock()
		}
		if !consensusValue.isNull{
			parseCommand(consensusValue.CommandData.Op, StateMachine)
		}
		c.Println("SM in seq",seq,":", StateMachine)	
		seq++
		
		//delete data to save memory

		/*
		for _,v:= range PQ {
			fmt.Print(*v)
		}
		fmt.Println()
		*/
		//deleteData(seq, 0)
		//time.Sleep(time.Duration(interval) * time.Millisecond)
	}

}

func weakMVC(stateStruct StateValueData , selfPort int, portNums []int,  seq int) (TerminationValue){

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
