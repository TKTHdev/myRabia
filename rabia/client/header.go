package main

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
)


type Data interface{}

type ConsensusData struct {
	Data Data
}

type CommandData struct {
	Op        string
	Timestamp int
	Seq       int
}

type Request struct {
	CommandData CommandData
	Redirected  bool
	Timestamp   int
}

type ResponseToClient struct {
	Value int
}

func init() {
	gob.Register(Request{})
	gob.Register(ResponseToClient{})
	gob.Register(ConsensusData{})
	gob.Register(CommandData{})

}

func receiveData(conn net.Conn) (ConsensusData, error) {
	var data ConsensusData
	fmt.Println("OKOK")
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("データ受信エラー:", err)
		return ConsensusData{}, err
	}
	fmt.Println("データ受信完了")
	return data, nil
}

func sendData(conn net.Conn, data Data) {
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(ConsensusData{Data: data})
	if err != nil {
		fmt.Println("データ送信エラー:", err)
		return
	}
}

func generateRandomCommand(readRatio int) string {
	variables := []string{"x", "y", "z", "a", "b", "c"}
	operators := []string{"=", "+", "-", "*", "%"}

	//generate read command with readRatio
	if rand.Intn(100) < readRatio {
		variable := variables[rand.Intn(len(variables))]
		return fmt.Sprintf("R %s", variable)
	}else{
		variable := variables[rand.Intn(len(variables))]
		operator := operators[rand.Intn(len(operators))]
		value := rand.Intn(100)
		return fmt.Sprintf("%s %s %d", variable, operator, value)
	}
}