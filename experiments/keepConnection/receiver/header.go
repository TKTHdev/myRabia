package main


import(
	"fmt"
	"net"
	"encoding/gob"
)

type Command struct{
	Op string
	Timestamp int

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