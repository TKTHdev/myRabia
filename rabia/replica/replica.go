package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
	"log"
	"os"
)

var conn net.Conn
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

	// プロキシからの接続を待ち受ける
	// 他のレプリカのポート番号を取得
	var portNums = listenAndAcceptConnectionWithProxy(listener, port)

	// 他のレプリカとの接続を確立
	go listenAndAccept(port)	
	time.Sleep(250 * time.Millisecond)


	//ここで合意アルゴリズムを実行
	var seq int = 0
	for i:=0;i<100;i++{
		fmt.Println("cnt: ", seq)

		command:=CommandData{Op: Operation, Timestamp: 0, Seq: seq}
		var stateValue int = exchangeStage(command, portNums,  seq)


		consensusValue :=weakMVC(stateValue,command, portInt, portNums, listener, seq)
		logger.Println("consensusValue: ", consensusValue)
		seq++
		
		//delete data to save memory
		if(seq>=2){
			deleteData(seq-2)
		}

		//time.Sleep(500 * time.Millisecond)
	}

}

func weakMVC(stateValue int ,command CommandData, selfPort int, portNums []int, ln net.Listener, seq int) int {

	var phase int = 0

	//Round 1
	var state StateValueData = StateValueData{Value: stateValue, Seq: seq, Phase: phase}
	voteValue :=roundOne(state, portNums,  seq)
	

	//Round 2
	var vote VoteValueData = VoteValueData{Value: voteValue, Seq: seq, Phase: phase}
	consensus,consensusValue :=roundTwo(vote, portNums, selfPort, seq,phase)
	if(consensus == -1){
		fmt.Println("reached consensus: ", consensusValue)
		return consensusValue
	}

	for{
		phase++
		state = StateValueData{Value: consensusValue, Seq: seq, Phase: phase}
		voteValue =roundOne(state, portNums,  seq)
		vote = VoteValueData{Value: voteValue, Seq: seq, Phase: phase}
		consensus,consensusValue =roundTwo(vote, portNums, selfPort, seq,phase)
		if(consensus == -1){
			fmt.Println("reached consensus: ", consensusValue)
			return  consensusValue
		}
	}

}
