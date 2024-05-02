package main 

import(
	"net"
	"fmt"
	"strconv"

)

var StateMachine map[string]int = make(map[string]int)

func main(){
	timestamp := 0
	portNumList := []int{8081, 8082, 8083}
	var choose string 
	fmt.Println("[A] to automatically send commands. [M] to manually send commands: ")
	fmt.Scan(&choose)
	if choose == "M"{
		for{
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
	}else if choose == "A"{
		var n int 
		fmt.Println("Enter the number of commands to send: ")
		fmt.Scan(&n)
		for i:=0; i<n; i++{
			var command string = generateRandomCommand()
			
			conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(portNumList[i%3]))
			if err != nil {
				fmt.Println("Dial error", err)
				return
			}
			defer conn.Close()
			sendData(conn, Request{CommandData: CommandData{Op: command, Timestamp: timestamp, Seq: 0}, Redirected: false, Timestamp: 0})
			parseCommand(command, StateMachine)
			timestamp++
		}
	}
	fmt.Println("SM: ", StateMachine)
}