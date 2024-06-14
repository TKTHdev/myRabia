package main

import (
	"fmt"
	"net"
)

var StateMachine map[string]int = make(map[string]int)

func main() {
	timestamp := 0
	IPList := []string{"52.63.13.55","13.237.225.199","54.253.27.126"};
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

				conn, err := net.Dial("tcp", IPList[i%3]+":8080")
				if err != nil {
					fmt.Println("Dial error", err)
					continue
				}
				defer conn.Close()
				sendData(conn, Request{CommandData: CommandData{Op: command, Timestamp: timestamp, Seq: 0}, Redirected: false, Timestamp: 0})
				if command[0] == 'R' {
					var data ConsensusData
					data, err :=receiveData(conn)
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
				}
		}
				timestamp++
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
				sendData(conn, Request{CommandData: CommandData{Op: command, Timestamp: timestamp, Seq: 0}, Redirected: false, Timestamp: 0})
				timestamp++
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
				sendData(conn, Request{CommandData: CommandData{Op: command, Timestamp: timestamp, Seq: 0}, Redirected: false, Timestamp: 0})
				timestamp++
		}
	}else{
		fmt.Println("Invalid input")
	}
}