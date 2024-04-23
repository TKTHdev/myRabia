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

func init() {
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
	portInt, _ := strconv.Atoi(port)

	// プロキシからの接続を待ち受ける
	// 他のレプリカのポート番号を取得
	var portNums, listener = listenAndAcceptConnectionWithProxy(listener, port)

	command := Command{Op: "write", Timestamp: 0}
	// 他のレプリカとの同期処理を実装
	for {
		weakMVC(command, portInt, portNums, listener)
	}

}

func weakMVC(command Command, selfPort int, portNums []int, listener net.Listener) {
	exchangeStage(command, portNums, listener, selfPort)
	time.Sleep(1000 * time.Millisecond)
}


func exchangeStage(command Command, portNums []int, listener net.Listener, port int) {
    // 他のレプリカとの同期処理を実装
    var state int
    wg := sync.WaitGroup{}
    wg.Add(1)
    go func() {
        state = exchangeBefore(command, portNums, port, &wg)
        wg.Done()
    }()
    go exchangeAfter(portNums, command, &wg)
    wg.Wait()
    fmt.Println("state: ", state)
}

func exchangeBefore(command Command, portNums []int, port int, wg *sync.WaitGroup) int {
    ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
    if err != nil {
        fmt.Println("リッスンエラー:", err)
        return -1
    }
    defer ln.Close()

    var receiveCnt int = 0
    var sameCnt int = 0

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("接続エラー:", err)
            return -1
        }

        go func(conn net.Conn) {
            defer conn.Close()
            receivedCommand, err := receiveCommand(conn)
            if err != nil {
                fmt.Println("コマンド受信エラー:", err)
                return
            }
            receiveCnt++
            if receivedCommand == command {
                sameCnt++
            }
        }(conn)

        if receiveCnt == len(portNums)/2+1 {
            if sameCnt == len(portNums)/2+1 {
				conn.Close()
                return 1
            }
			conn.Close()
            return 0
        }
		time.Sleep(1 * time.Second)
    }
}

func exchangeAfter(portNums []int, command Command, wg *sync.WaitGroup) {
    for _, portNum := range portNums {
        conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(portNum))
        if err != nil {
			continue	
        }
		fmt.Println("sousinn")
        defer conn.Close()
        sendCommand(conn, command)
    }
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
	defer func() {
		if err := listener.Close(); err != nil {
			fmt.Println("クローズエラー:", err)
		}
	}()

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
