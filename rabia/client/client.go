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
	fmt.Println("[A] to automatically send commands. [M] to manually send commands: ")
	fmt.Scan(&choose)
	for {

		if choose == "M" {
			for {
				var command string
				fmt.Println("Enter the command: ")
				fmt.Scan(&command)

				conn, err := net.Dial("tcp", "localhost:8081")
				if err != nil {
					fmt.Println("Dial error", err)
					return
				}
				defer conn.Close()
				fmt.Println("Connected to the server")
				sendData(conn, Request{CommandData: CommandData{Op: command, Timestamp: timestamp, Seq: 0}, Redirected: false, Timestamp: 0})
				timestamp++

			}
		} else if choose == "A" {
			var n int
			fmt.Println("Enter the number of commands to send: ")
			fmt.Scan(&n)
			for i := 0; i < n; i++ {
				var command string = generateRandomCommand()

				conn, err := net.Dial("tcp", IPList[i%3]+":8080")
				if err != nil {
					fmt.Println("Dial error", err)
					continue
				}
				defer conn.Close()
				sendData(conn, Request{CommandData: CommandData{Op: command, Timestamp: timestamp, Seq: 0}, Redirected: false, Timestamp: 0})
				parseCommand(command, StateMachine)
				timestamp++
			}
		}
		fmt.Println("SM: ", StateMachine)
	}
}
