package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func sendCommand(conn net.Conn, command Command) {
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(command)
	if err != nil {
		return
	}
}

func receiveCommand(conn net.Conn) (Command, error) {
	var data Command
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&data)
	if err != nil {

		return Command{}, err
	}
	return data, nil
}

func sendBinaryValue(conn net.Conn, value int) {
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(value)
	if err != nil {
		return
	}
}

func receiveBinaryValue(conn net.Conn) (int, error) {
	var data int
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&data)
	if err != nil {
		return -1, err
	}
	return data, nil
}

type Command struct {
	Op        string
	Timestamp int
}

func setConnectionWithOtherReplicas(portNums []int, selfPort int) []net.Conn {
	var conns []net.Conn

	for _, portNum := range portNums {
		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", portNum))
		if err != nil {
			fmt.Printf("接続エラー (ポート番号: %d): %v\n", portNum, err)
			continue
		}

		conns = append(conns, conn)
	}
	fmt.Println(conns)
	return conns
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

func listenAndAcceptConnectionWithProxy(listener net.Listener, port string) []int {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("リッスンエラー:", err)
		return nil
	}
	defer func() {
		if err := listener.Close(); err != nil {
			fmt.Println("クローズエラー:", err)
		}
	}()

	// リクエストを待ち受ける
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("接続エラー:", err)
		return nil
	}
	fmt.Println("接続しました")

	//プロキシからポート番号リストを読み取る
	var message string
	for {
		messagePart, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("メッセージ読み取りエラー:", err)
			conn.Close()
			return nil
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
		return nil
	}
	fmt.Println("他のレプリカのポート番号リスト: ", portNums)

	return portNums
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
