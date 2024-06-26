package main

import (
	"fmt"
	"net"
	"time"
)

var StateMachine map[string]int = make(map[string]int)
var IPList = []string{"52.62.115.28", "52.65.112.127", "52.64.108.149"}
var replicaNum = 3

func main() {
	command , clientNum, duration := setUp()
	var stopChannelList []chan bool = make([]chan bool, clientNum)
	var commandNumChannelList []chan int = make([]chan int, clientNum)

	for i := 0; i < clientNum; i++ {
		stopChannelList[i] = make(chan bool)
		commandNumChannelList[i] = make(chan int)
	}

	for i := 0; i < clientNum; i++ {
		  go YCSB(command, stopChannelList[i],commandNumChannelList[i] , i)
	}
	time.Sleep(time.Duration(duration) * time.Second)

	fmt.Println("Test finished")

	//stop each of goroutines
	for i := 0; i < clientNum; i++ {
		close(stopChannelList[i])
	}

	//get the number of commands executed by each client
	var totalCommandNum int = 0
	for i := 0; i < clientNum; i++ {
		fmt.Println("client", i ,"stop")
		totalCommandNum += <-commandNumChannelList[i]
	}

	fmt.Println("Total number of commands executed: ", totalCommandNum)


}


func YCSB(command string, stopChannel chan bool, commandNumChannel chan int,  ID int)  {

	var readRatio int	

	if command == "A" {
		readRatio = 50
	} else if command == "B" {
		readRatio = 95
	} else if command == "C" {
		readRatio = 100
	}


	var cnt int = 0
	
	var replicaID string = IPList[ID%replicaNum]+":8080"

	for  {
			select{
			case <-stopChannel:
				fmt.Println("Client stopped")
				return 

			default:
			var command string = generateRandomCommand(readRatio)
			//fmt.Println("Command: " + command)
			conn, err := net.Dial("tcp", replicaID)
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
						//fmt.Println("Key not found")
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
			cnt++
		}
	}
}


func setUp()(string, int , int){
	fmt.Println("YCSB [A] or [B] or [C] ?")
	var command string
	fmt.Scanln(&command)
	fmt.Println("Number of concurrent clients?")
	var clientNum int
	fmt.Scanln(&clientNum)
	fmt.Println("How long to run the test for (in seconds)?")
	var time int
	fmt.Scanln(&time)
	return command, clientNum, time
}