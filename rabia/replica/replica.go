package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"time"
)

var conn net.Conn
var listener net.Listener

func init() {
	gob.Register(Command{})
}

func main() {
	// サーバーに接続し、自身に割り当てられたポート番号を受け取る
	var Operation string
	fmt.Println("Operation: ")
	fmt.Scan(&Operation)
	command := Command{Op: Operation, Timestamp: 0}
	var port string = sayHelloAndReceivePortNum(conn)
	portInt, _ := strconv.Atoi(port)

	// プロキシからの接続を待ち受ける
	// 他のレプリカのポート番号を取得
	var portNums = listenAndAcceptConnectionWithProxy(listener, port)

	ln, err := net.Listen("tcp", ":"+strconv.Itoa(portInt))
	if err != nil {
		fmt.Println("リッスンエラー:", err)
		return
	}

	time.Sleep(250 * time.Millisecond)

	// 他のレプリカとの同期処理を実装
	weakMVC(command, portInt, portNums, ln)

}

func weakMVC(command Command, selfPort int, portNums []int, ln net.Listener) {
	// Exchange Stage
	var state int
	state = exchangeStage(command, portNums, selfPort, ln)
	fmt.Println("state: ", state)
	var vote int = roundOne(state, portNums, selfPort, ln)
	fmt.Println("vote: ", vote)
}
