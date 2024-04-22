package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)



var conn net.Conn
var listener net.Listener
var wg sync.WaitGroup



func init(){
	gob.Register(Command{})
}

func sendCommand(conn net.Conn, command Command) {
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(command)
	if err != nil {
		fmt.Println("エンコードエラー:", err)
		return
	}
}

func receiveCommand(conn net.Conn) (Command, error) {
	var data Command
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("デコードエラー:", err)
		return Command{}, err
	}
	return data, nil
}

func main() {
	// サーバーに接続し、自身に割り当てられたポート番号を受け取る
	var port string = sayHelloAndReceivePortNum(conn)

	// プロキシからの接続を待ち受ける
	// 他のレプリカのポート番号を取得
	var portNums, listener = listenAndAcceptConnectionWithProxy(listener, port)

	command:= Command{Op: "write", Timestamp: 0}
	// 他のレプリカとの同期処理を実装
	for {
		weakMVC(command,port, portNums, listener)
		seq++
	}
	
}

func weakMVC(command Command, selfPort string, portNums []int, listener net.Listener) {
	return 
}

func exchangeBefore(command Command,selfPort string, portNums []int, listener net.Listener) {
	var cnt int = 0

	for _, portNum := range portNums {
		go func(portNum int) {
			
		}()
	}	

}


func exchangeAfter(port string, portNums []int, listener net.Listener) {

}



func sayHelloAndReceivePortNum(conn net.Conn) string {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("接続エラー:", err)
		return ""
	}
	defer conn.Close()

	message := "Hello\n"
	_, err = fmt.Fprint(conn, message)
	if err != nil {
		fmt.Println("メッセージ送信エラー:", err)
		return ""
	}
	fmt.Println("メッセージを送信しました")

	// サーバーから、自身に割り当てられたポート番号を受け取る
	reply, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("ポート番号受信エラー:", err)
		return ""
	}
	fmt.Print("サーバーからの返信(割り当てられたポート番号): ", reply)

	// ポート番号を取得
	reply = strings.TrimSpace(reply)
	return reply
}

func listenAndAcceptConnectionWithProxy(listener net.Listener, port string) ([]int, net.Listener) {
	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		fmt.Println("リッスンエラー:", err)
		return nil, nil
	}

	// リクエストを待ち受ける
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("接続エラー:", err)
		return nil, nil
	}
	fmt.Println("接続しました")

	//プロキシからポート番号リストを読み取る
	var message string
	for {
		messagePart, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("メッセージ読み取りエラー:", err)
			conn.Close()
			return nil, nil
		}
		message += messagePart
		if strings.HasSuffix(message, "\n") {
			break
		}
	}
	//クライアントからのメッセージを表示
	fmt.Print("クライアントからのメッセージ: ", message)

	time.Sleep(1000 * time.Millisecond)
	//ポート番号リストをパース
	portNums, err := parsePortList(message)
	if portNums == nil {
		return nil, nil
	}
	fmt.Println("他のレプリカのポート番号リスト: ", portNums)

	return portNums, listener
}

func parsePortList(portList string) ([]int, error) {
	var portNums []int

	// スペースでポート番号文字列を分割
	portStrings := strings.Fields(portList)

	// 各ポート番号文字列を整数に変換
	for _, portString := range portStrings {
		portNum, err := strconv.Atoi(portString)
		if err != nil {
			return nil, fmt.Errorf("無効なポート番号: %s", portString)
		}
		portNums = append(portNums, portNum)
	}

	return portNums, nil
}

