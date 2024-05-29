package main

import (
	"fmt"
	"net"
	"sync"
)

var wg sync.WaitGroup
var replicaIPs []string
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
		go handleConnection(conn)
		cnt++
		wg.Wait()
		if cnt == n {
			fmt.Println("All replicas connected")
			break
		}
	}
}

func handleConnection(conn net.Conn) {
	//Calculate the port number of the replica
	//レプリカのポート番号を計算
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("クローズエラー:", err)
		}
	}(conn)
	fmt.Println("Connected to replica: ", conn.RemoteAddr().String())

	replicaIPs = append(replicaIPs, conn.RemoteAddr().String())
}

// Send the list of port numbers to the replicas
// レプリカにポート番号のリストを送信
func sendPortNumListToReplicas() {
	for _, IP := range replicaIPs {
		wg.Add(1)
		go func(IP string) {
			conn, err := net.Dial("tcp", IP+":8080")
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
			_, err = fmt.Fprintf(conn, portListToString()+"\n")
			if err != nil {
				return
			}
			fmt.Printf("Port number list sent to replica %d\n", net.IPv4len)
			wg.Done()
		}(IP)
		wg.Wait()
	}
}

// Convert the list of port numbers to a string
// ポート番号のリストを文字列に変換
func portListToString() string {
	var IPListString string
	for i, IP := range replicaIPs {
		IPListString += IP
		if i < len(replicaIPs)-1 {
			IPListString += ","
		}
	}
	return IPListString
}
