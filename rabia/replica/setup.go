package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func setConnectionWithOtherReplicas(portNums []int) []net.Conn {
	var conns []net.Conn

	for _, portNum := range portNums {
		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", portNum))
		if err != nil {
			//fmt.Printf("接続エラー (ポート番号: %d): %v\n", portNum, err)
			continue
		}

		conns = append(conns, conn)
	}
	//fmt.Println(conns)
	return conns
}

func RegisterToProxy() {
	conn, err := net.Dial("tcp", "13.236.12.56:8080")
	if err != nil {
		fmt.Println("接続エラー:", err)
		return
	}
	defer conn.Close()

	fmt.Println("プロキシに接続しました")
	return
}

func listenAndAcceptConnectionWithProxy() []string {
	listener, err := net.Listen("tcp", ":8080")
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
	fmt.Println("レプリカから返事を受けました")

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

	time.Sleep(1000 * time.Millisecond)
	//ポート番号リストをパース
	IPs := strings.Split(message, ",")
	if IPs == nil {
		return nil
	}
	fmt.Println("他のレプリカのIPリスト: ", IPs)

	return IPs
}

func parseIPList(IPList string) ([]string, error) {
	return strings.Split(IPList, ","), nil
}
