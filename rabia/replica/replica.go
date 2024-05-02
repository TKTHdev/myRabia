package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var listener net.Listener

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
	var Operation string
	fmt.Println("Operation: ")
	fmt.Scan(&Operation)
	var port string = sayHelloAndReceivePortNum()
	portInt, _ := strconv.Atoi(port)
	command := CommandData{Op: Operation, Timestamp: 0, Seq: 0}

	// プロキシからの接続を待ち受ける
	// 他のレプリカのポート番号を取得
	var portNums = listenAndAcceptConnectionWithProxy(listener, port)

	// 他のレプリカとの接続を確立
	go listenAndAccept(port)	
	time.Sleep(250 * time.Millisecond)


	//ここで合意アルゴリズムを実行
	var seq int = 0
	for {
		fmt.Println("cnt: ", seq)

		PQMutex.Lock()
		commandPointer :=PQ.Head()
		if commandPointer == nil{
			PQMutex.Unlock()
			continue
		}
		PQMutex.Unlock()
		var terminationFlag,stateValue int = exchangeStage(command, portNums,  seq)
		if terminationFlag == 1{
			logger.Println("consensusValue: ", stateValue)
			seq++
			continue
		}

	    consensusValue:=weakMVC(stateValue,command, portInt, portNums, listener, seq)
		logger.Println("consensusValue: ", consensusValue)
		seq++
		
		//delete data to save memory
		

		//time.Sleep(500 * time.Millisecond)

	}

}

func weakMVC(stateValue int ,command CommandData, selfPort int, portNums []int, ln net.Listener, seq int) int {

	var phase int = 0

	//Round 1
	var state StateValueData = StateValueData{Value: stateValue, Seq: seq, Phase: phase}
	terminationFlag, voteValue :=roundOne(state, portNums,  seq, phase)
	if terminationFlag == 1{
		fmt.Println("reached consensus: ", voteValue)
		return voteValue
	}	
	
	
	//Round 2
	var vote VoteValueData = VoteValueData{Value: voteValue, Seq: seq, Phase: phase}
	terminationFlag,consensusValue :=roundTwo(vote, portNums, selfPort, seq,phase)
	if(terminationFlag == 1){
		fmt.Println("reached consensus: ", consensusValue)
		return consensusValue 
	}

	for{
		phase++

		state = StateValueData{Value: consensusValue, Seq: seq, Phase: phase}
		terminationFlag, voteValue =roundOne(state, portNums,  seq, phase)
		if terminationFlag == 1{
			fmt.Println("reached consensus: ", voteValue)
			return voteValue
		}

		vote = VoteValueData{Value: voteValue, Seq: seq, Phase: phase}
		terminationFlag,consensusValue =roundTwo(vote, portNums, selfPort, seq,phase)
		if terminationFlag == 1{
			fmt.Println("reached consensus: ", consensusValue)
			return  consensusValue
		}

		deleteData(seq, phase)

	}

}
