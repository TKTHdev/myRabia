package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"sync"
)

var wg sync.WaitGroup
var portNums []int
var listener net.Listener

func main() {
	//Set the number of replicas
	//レプリカの数を決める
	var n int
	fmt.Println("Enter the number of replicas: ")
	_, err := fmt.Scan(&n)
	if err != nil {
		return
	}

	//Create a slice of connections
	//サーバーを起動
	listener = initServer(listener)
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("クローズエラー:", err)
		}
	}(listener)
	//Accept connections from n replicas
	//n個のレプリカからの接続を受け付ける
	//Wait for n replicas to connect
	//n個のレプリカが接続するのを待つ
	acceptNConnections(listener, n)

	//send the list of port numbers to the replicas
	//レプリカにポート番号のリストを送信
	sendPortNumListToReplicas()
}

func initServer(listener net.Listener) net.Listener {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("リッスンエラー:", err)
		return nil
	}

	fmt.Println("サーバーが起動しました。クライアントからの接続を待機しています...")
	return listener
}

func acceptNConnections(listener net.Listener, n int) {
	var cnt = 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("接続エラー:", err)
			continue
		}
		fmt.Println("接続しました")
		wg.Add(1)
		go handleConnection(conn, cnt)
		cnt++
		wg.Wait()
		if cnt == n {
			fmt.Println("All replicas connected")
			break
		}
	}
}

func handleConnection(conn net.Conn, portNumOffset int) {
	//Calculate the port number of the replica
	//レプリカのポート番号を計算
	var portNum = 8081 + portNumOffset

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("クローズエラー:", err)
		}
	}(conn)

	//Read the message from the replica
	//レプリカからのメッセージを読み取る
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("メッセージ読み取りエラー:", err)
		return
	}
	fmt.Printf("Message from replica %d\n", portNumOffset+1)

	if message == "q\n" {
		fmt.Println("クライアントが接続を終了しました")
		return
	}
	//Send the port number to the replica
	//レプリカにポート番号を送信
	sendPortNumToReplica(conn, portNum)
	portNums = append(portNums, portNum)
}

// Send the port number to the replica
// レプリカにポート番号を送信
func sendPortNumToReplica(conn net.Conn, portNum int) {
	fmt.Println("Sending port number to replica")
	_, i := fmt.Fprintf(conn, strconv.Itoa(portNum)+"\n")
	if i != nil {
		return
	}
	fmt.Printf("Port number %d sent to replica\n", portNum)
	wg.Done()
}

// Send the list of port numbers to the replicas
// レプリカにポート番号のリストを送信
func sendPortNumListToReplicas() {
	for _, portNum := range portNums {
		wg.Add(1)
		go func(portNum int) {
			fmt.Print("Connecting to replica ", portNum, "...")
			conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(portNum))
			if err != nil {
				fmt.Println("接続エラー:", err)
				return
			}
			defer func(conn net.Conn) {
				err := conn.Close()
				if err != nil {
					fmt.Println("クローズエラー:", err)
				}
			}(conn)
			fmt.Println("Sending port number list to replica")
			_, err = fmt.Fprintf(conn, portListToString(portNums)+"\n")
			if err != nil {
				return
			}
			fmt.Printf("Port number list sent to replica %d\n", portNum)
			wg.Done()
		}(portNum)
		wg.Wait()
	}
}

// Convert the list of port numbers to a string
// ポート番号のリストを文字列に変換
func portListToString(portNums []int) string {
	var portList string
	for _, portNum := range portNums {
		portList += strconv.Itoa(portNum) + " "
	}
	return portList
}
