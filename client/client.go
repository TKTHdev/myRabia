package main

import (
	"fmt"
	"net"
	"time"
)

var StateMachine map[string]int = make(map[string]int)
var IPList = []string{"52.62.115.28", "52.65.112.127", "52.64.108.149"}

func main() {
	fmt.Println("YCSB [A] or [B] or [C] ?")
	var command string
	fmt.Scanln(&command)
	fmt.Println("Number of commands?")
	var commandNum int
	fmt.Scanln(&commandNum)
	YCSB(command, commandNum)
}

func YCSB(command string, commandNum int){

	var readRatio int	


	if command == "A" {
		readRatio = 50
	} else if command == "B" {
		readRatio = 95
	} else if command == "C" {
		readRatio = 100
	}

	for i := 0; i < commandNum; i++ {
		var command string = generateRandomCommand(readRatio)
		fmt.Println("Command: " + command)
		conn, err := net.Dial("tcp", IPList[i%3]+":8080")
		if err != nil {
			fmt.Println("Dial error", err)
			continue
		}
		defer conn.Close()
		sendData(conn, Request{CommandData: CommandData{Op: command, Timestamp: time.Now(), Seq: 0, ClientAddr: conn.LocalAddr().String()}, Redirected: false, Timestamp: 0})
		if command[0] == 'R' {
			var data ConsensusData
			data, err := receiveData(conn)
			if err != nil {
				fmt.Println("Error in receiving data")
			}
			response := data.Data
			switch response := response.(type) {
			case ResponseToClient:
				if response.Value == -1 {
					fmt.Println("Key not found")
				}
				if response.Value != -1 {
					//fmt.Println("Read value: ", response.Value)
				}
			}
		} else {
			var data ConsensusData
			data, err := receiveData(conn)
			if err != nil {
				fmt.Println("Error in receiving data")
			}
			response := data.Data
			switch response := response.(type) {
			case ResponseToClient:
				if response.Value == 0 {
					//fmt.Println("Write successful")
				} else {
					fmt.Println("Write unsuccessful")
				}
			}
		}
		fmt.Println("Command ", i+1, " completed")
	}
}