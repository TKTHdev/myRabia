package main

import (
	"fmt"
	"net"
	"time"
)

var StateMachine map[string]int = make(map[string]int)
var IPList = []string{"52.64.108.149", "52.65.112.127", "52.62.115.28"}
var replicaNum = 3

type Report struct {
	commandNum int
	readTime   time.Duration
	writeTime  time.Duration
	totalTime  time.Duration
}

func main() {
	command, clientNum, duration := setUp()
	var stopChannelList []chan bool = make([]chan bool, clientNum)
	var reportChannelList []chan Report = make([]chan Report, clientNum)

	for i := 0; i < clientNum; i++ {
		stopChannelList[i] = make(chan bool)
		reportChannelList[i] = make(chan Report)
	}

	for i := 0; i < clientNum; i++ {
		go YCSB(command, stopChannelList[i], reportChannelList[i], i)
	}
	time.Sleep(time.Duration(duration) * time.Second)

	fmt.Println("Test finished")

	//stop each of goroutines
	for i := 0; i < clientNum; i++ {
		close(stopChannelList[i])
	}

	//get the number of commands executed by each client
	var totalCommandNum int = 0
	var totalReadTime time.Duration = 0
	var totalWriteTime time.Duration = 0
	var totalTime time.Duration = 0

	for i := 0; i < clientNum; i++ {
		fmt.Println("client", i, "stop")
		report := <-reportChannelList[i]
		totalCommandNum += report.commandNum
		totalReadTime += report.readTime
		totalWriteTime += report.writeTime
		totalTime += report.totalTime
	}

	fmt.Println("Total number of commands executed: ", totalCommandNum)
	fmt.Printf("Average read time: %d ms\n", totalReadTime.Milliseconds()/int64(clientNum))
	fmt.Printf("Average write time: %d ms\n", totalWriteTime.Milliseconds()/int64(clientNum))
	fmt.Printf("Average total time: %d ms\n", totalTime.Milliseconds()/int64(clientNum))
}

func YCSB(command string, stopChannel chan bool, reportChannel chan Report, ID int) {
	var readRatio int

	if command == "A" {
		readRatio = 50
	} else if command == "B" {
		readRatio = 95
	} else if command == "C" {
		readRatio = 100
	}

	var cnt int = 0

	var replicaID string = IPList[ID%replicaNum] + ":8080"
	//fmt.Println("Replica ID: ", replicaID)
	conn, err := net.Dial("tcp", replicaID)
	if err != nil {
		fmt.Println("Dial error", err)
		return
	}

	var readTime time.Duration = 0
	var writeTime time.Duration = 0
	var total time.Duration = 0
	var readCnt int = 0
	var writeCnt int = 0

	for {
		select {
		case <-stopChannel:
			readTimeAverage := readTime / time.Duration(readCnt)
			writeTimeAverage := writeTime / time.Duration(writeCnt)
			fmt.Println("Client stopped")
			reportChannel <- Report{commandNum: cnt, readTime: readTimeAverage, writeTime: writeTimeAverage, totalTime: total}
			return

		default:
			var command string = generateRandomCommand(readRatio)
			//fmt.Println("Command: " + command)
			start := time.Now()
			sendData(conn, Request{CommandData: CommandData{Op: command, Timestamp: time.Now(), Seq: 0, ClientAddr: conn.LocalAddr().String()}, Redirected: false, Timestamp: 0})

			if command[0] == 'R' {
				//start measuring time
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
				//end measuring time
				elapsed := time.Since(start)
				readTime += elapsed
				total += elapsed
				readCnt++
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

				elapsed := time.Since(start)
				writeTime += elapsed
				total += elapsed
				writeCnt++
			}
			cnt++
		}
	}
}

func setUp() (string, int, int) {
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