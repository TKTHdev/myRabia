package main

import (
	"container/heap"
	"encoding/gob"
	"fmt"
	"net"
	"sync"
)

var CommandDataMutex sync.Mutex
var StateValueDataMutex sync.Mutex
var VoteValueDataMutex sync.Mutex
var ConsensusTerminationMutex sync.Mutex
var PQMutex sync.Mutex

type Data interface{}

type ConsensusData struct {
	Data Data
}

type CommandData struct {
	Op        string
	Timestamp int
	Seq       int
}

type StateValueData struct {
	Value       int
	Seq         int
	Phase       int
	CommandData CommandData
}

type VoteValueData struct {
	Value       int
	Seq         int
	Phase       int
	CommandData CommandData
}

type SeqPhase struct {
	Seq   int
	Phase int
}

type RoundTwoReturnStruct struct {
	ConsensusValue int
	CommandData    CommandData
}

type ConsensusTermination struct {
	Seq         int
	Value       int
	CommandData CommandData
}

type TerminationValue struct {
	isNull      bool
	CommandData CommandData
	phase       int
	seq         int
}

type Request struct {
	CommandData CommandData
	Redirected  bool
	Timestamp   int
}

type CommandTimestamp struct {
	Command   string
	Timestamp int
}

type ResponseToClient struct {
	Value int
}

var CommandDataMapList map[int][]CommandData
var StateValueDataMapList map[SeqPhase][]StateValueData
var VoteValueDataMapList map[SeqPhase][]VoteValueData
var ConsensusTerminationMapList map[int][]ConsensusTermination
var Dictionary map[CommandTimestamp]bool
var PQ PriorityQueue

func init() {
	CommandDataMapList = make(map[int][]CommandData)
	StateValueDataMapList = make(map[SeqPhase][]StateValueData)
	VoteValueDataMapList = make(map[SeqPhase][]VoteValueData)
	ConsensusTerminationMapList = make(map[int][]ConsensusTermination)
	Dictionary = make(map[CommandTimestamp]bool)

	gob.Register(CommandData{})
	gob.Register(StateValueData{})
	gob.Register(VoteValueData{})
	gob.Register(ConsensusTermination{})
	gob.Register(Request{})
	gob.Register(ResponseToClient{})
	gob.Register(ConsensusData{})

	PQ := make(PriorityQueue, 0)
	heap.Init(&PQ)

}

func listenAndAccept() {
	ln, err := net.Listen("tcp", ":8080")
	//fmt.Println("Listening on port: ", port)
	if err != nil {
		fmt.Println("リッスンエラー:", err)
		return
	}
	defer ln.Close()

	for {
		//fmt.Println("Waiting for connection...")
		conn, err := ln.Accept()
		//fmt.Println("Message received from: ", conn.RemoteAddr())
		if err != nil {
			//fmt.Println("接続エラー:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func sendData(conn net.Conn, data Data) {
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(ConsensusData{Data: data})
	if err != nil {
		fmt.Println("データ送信エラー:", err)
		return
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		consensusData, err := receiveData(conn)
		if err != nil {
			//fmt.Println("データ受信エラー:", err)
			return
		}
		data := consensusData.Data

		switch data := data.(type) {
		case CommandData:
			CommandDataMutex.Lock()
			//fmt.Println("Received CommandData: ", data)
			if data.Op[0]!='R'{
				CommandDataMapList[data.Seq] = append(CommandDataMapList[data.Seq], data)
			}else{
				value, err:= parseReadCommand(data.Op, StateMachine)
				if err == "notFound"{
					fmt.Println("Addresses of client and server: ", conn.RemoteAddr(), conn.LocalAddr())
					response :=  ResponseToClient{Value: -1}
					sendData(conn, response)
				}else{
					response := ResponseToClient{Value: value}
					sendData(conn, response)
				}
				fmt.Println("Sent")
			}
			//fmt.Println("CommandDataMapList: ", CommandDataMapList)
			CommandDataMutex.Unlock()
		case StateValueData:
			StateValueDataMutex.Lock()
			//fmt.Println("Received StateValueData: ", data)
			StateValueDataMapList[SeqPhase{Seq: data.Seq, Phase: data.Phase}] = append(StateValueDataMapList[SeqPhase{Seq: data.Seq, Phase: data.Phase}], data)
			//fmt.Println("StateValueDataMapList: ", StateValueDataMapList)
			StateValueDataMutex.Unlock()
		case VoteValueData:
			VoteValueDataMutex.Lock()
			//fmt.Println("Received VoteValueData: ", data)
			VoteValueDataMapList[SeqPhase{Seq: data.Seq, Phase: data.Phase}] = append(VoteValueDataMapList[SeqPhase{Seq: data.Seq, Phase: data.Phase}], data)
			//fmt.Println("VoteValueDataMaplist: ", VoteValueDataMaplist)
			VoteValueDataMutex.Unlock()
		case ConsensusTermination:
			ConsensusTerminationMutex.Lock()
			//fmt.Println("Received ConsensusTermination: ", data)
			ConsensusTerminationMapList[data.Seq] = append(ConsensusTerminationMapList[data.Seq], data)
			ConsensusTerminationMutex.Unlock()

		case Request:
			PQMutex.Lock()
			//fmt.Println("Received Request: ", data)
			if !data.Redirected {
				data.Redirected = true
				broadCastData(replicaIPs, data)
			} else {
				PQ.Push(&data.CommandData)
			}
			PQMutex.Unlock()

		default:
			fmt.Println("未知のデータ型です:", data)
		}
	}
}

func broadCastData(IPLists []string, data Data) {
	conns := setConnectionWithOtherReplicas(IPLists)
	for _, conn := range conns {
		sendData(conn, data)
	}
}

func receiveData(conn net.Conn) (ConsensusData, error) {
	var data ConsensusData
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&data)
	if err != nil {
		return ConsensusData{}, err
	}
	return data, nil
}

func deleteData(seq int, phase int) {
	CommandDataMutex.Lock()
	delete(CommandDataMapList, seq)
	CommandDataMutex.Unlock()

	StateValueDataMutex.Lock()
	delete(StateValueDataMapList, SeqPhase{Seq: seq, Phase: phase})
	StateValueDataMutex.Unlock()

	VoteValueDataMutex.Lock()
	delete(VoteValueDataMapList, SeqPhase{Seq: seq, Phase: phase})
	VoteValueDataMutex.Unlock()
}

//Unused now
//Measure the size of the map data structure

/*
func printCommandDataMapListSize() {
    CommandDataMutex.Lock()
    var totalSize int
    for _, value := range CommandDataMapList {
        size := len(value)
        totalSize += size
    }
    fmt.Printf("Total Size: %d\n", totalSize)
    CommandDataMutex.Unlock()
}

func printStateValueDataMapListSize() {
    StateValueDataMutex.Lock()
    var totalSize int
    for _, value := range StateValueDataMapList {
        size := len(value)
        totalSize += size
    }
    fmt.Printf("Total Size: %d\n", totalSize)
    StateValueDataMutex.Unlock()
}

func printVoteValueDataMapListSize() {
    VoteValueDataMutex.Lock()
    var totalSize int
    for _, value := range VoteValueDataMapList {
        size := len(value)
        totalSize += size
    }
    fmt.Printf("Total Size: %d\n", totalSize)
    VoteValueDataMutex.Unlock()
}
*/
