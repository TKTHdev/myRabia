package main

import (
	"fmt"
	"net"
	"encoding/gob"
)


type Command struct{
	Op string
	Timestamp int

}

func sendCommand(conn net.Conn, command Command) {
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(command)
	if err != nil {
		fmt.Println("エンコードエラー:", err)
		return
	}
}

