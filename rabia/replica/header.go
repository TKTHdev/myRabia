package main


import(
	"encoding/gob"
	"net"

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

type Command struct{
	Op string
	Timestamp int

}

