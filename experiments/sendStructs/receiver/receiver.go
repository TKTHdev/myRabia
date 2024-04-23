package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println("リッスンエラー:", err)
		return
	}
	defer ln.Close()
	for{
		conn,err:= ln.Accept()
		if err != nil {
			fmt.Println("接続エラー:", err)
			return
		}
		go func(conn net.Conn){
			defer conn.Close()
			data,err:= receiveCommand(conn)
			if err != nil {
				fmt.Println("コマンド受信エラー:", err)
				return
			}
			fmt.Println("コマンドを受信しました: ", data)
		}(conn)
	}
}
