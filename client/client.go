package main

import (
	"fmt"
	"net"
)

var StateMachine map[string]int = make(map[string]int)

func main() {
	timestamp := 0
	IPList := []string{"52.62.115.28", "52.65.112.127", "52.64.108.149"}
	var choose string
	var commandNum int
	fmt.Println("YCSB Workload [A] or [B] or [C] ?: ")
	fmt.Println("A: 50% Read, 50% Write")
	fmt.Println("B: 95% Read, 5% Write")
	fmt.Println("C: 100% Read")
	fmt.Scan(&choose)
	fmt.Println("How many commands do you want to run?: ")
	fmt.Scan(&commandNum)

	if choose == "A" {
		for i := 0; i < commandNum; i++ {
			var command string = generateRandomCommand(50)
			fmt.Println("Command: " + command)
			conn, err := net.Dial("tcp", IPList[i%3]+":8080")
			if err != nil {
				fmt.Println("Dial error", err)
				continue
			}
			defer conn.Close()
			sendData(conn, Request{CommandData: CommandData{Op: command, Timestamp: timestamp, Seq: 0, ClientAddr: conn.LocalAddr().String()}, Redirected: false, Timestamp: 0})
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

	if choose == "B" {
		for i := 0; i < commandNum; i++ {
			var command string = generateRandomCommand(95)

			conn, err := net.Dial("tcp", IPList[i%3]+":8080")
			if err != nil {
				fmt.Println("Dial error", err)
				continue
			}
			defer conn.Close()
			sendData(conn, Request{CommandData: CommandData{Op: command, Timestamp: timestamp, Seq: 0, ClientAddr: conn.LocalAddr().String()}, Redirected: false, Timestamp: 0})
			if command[0] == 'R' {
				var data ConsensusData
				fmt.Println("READ")
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
						fmt.Println("Read value: ", response.Value)
					}
				}
			} else {
				var data ConsensusData
				fmt.Println("WRITE")
				data, err := receiveData(conn)
				if err != nil {
					fmt.Println("Error in receiving data")
				}
				response := data.Data
				switch response := response.(type) {
				case ResponseToClient:
					if response.Value == 0 {
						fmt.Println("Write successful")
					} else {
						fmt.Println("Write unsuccessful")
					}
				}
			}
		}
	}

	if choose == "C" {
		for i := 0; i < commandNum; i++ {
			var command string = generateRandomCommand(100)

			conn, err := net.Dial("tcp", IPList[i%3]+":8080")
			if err != nil {
				fmt.Println("Dial error", err)
				continue
			}
			defer conn.Close()
			sendData(conn, Request{CommandData: CommandData{Op: command, Timestamp: timestamp, Seq: 0, ClientAddr: conn.LocalAddr().String()}, Redirected: false, Timestamp: 0})
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
						fmt.Println("Read value: ", response.Value)
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
						fmt.Println("Write successful")
					} else {
						fmt.Println("Write unsuccessful")
					}
				}
			}
		}
	}
}
