package main 

import(
	"net"
	"fmt"

)


func main(){
	timestamp := 0
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
}