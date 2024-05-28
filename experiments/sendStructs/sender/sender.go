package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	for {

		conn, err := net.Dial("tcp", "localhost:9999")
		if err != nil {
			continue
		}
		defer conn.Close()

		data := Command{"GET", 123}
		sendCommand(conn, data)
		fmt.Println("コマンドを送信しました: ", data)
		time.Sleep(30 * time.Microsecond)
	}
}
