package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func setConnectionWithOtherReplicas(IPLists []string) []net.Conn {
	var conns []net.Conn

	for _, IP := range IPLists {
		if IP[0]!= 'N'  && IP!= stringIP{
			conn, err := net.Dial("tcp", IP+":8080")
			if err != nil {
				//fmt.Printf("接続エラー  %d: %v\n", conn, err)
				continue
			}
			conns = append(conns, conn)
		}
	}
	//fmt.Println(conns)
	return conns
}

func RegisterToProxy() string {
	conn, err := net.Dial("tcp", "13.237.227.47:8080")
	if err != nil {
		fmt.Println("接続エラー:", err)
		return ""
	}
	defer conn.Close()

	fmt.Println("プロキシに接続しました")
	removePort := strings.Split(conn.LocalAddr().String(), ":")
	return removePort[0]
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
	//\n を削除
	message = strings.TrimRight(message, "\n")

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
