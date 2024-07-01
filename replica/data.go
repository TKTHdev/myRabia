package main

import (
	"container/heap"
	"encoding/gob"
	"fmt"
	"time"
	"net"
	"sync"
	"github.com/chyeh/pubip"
)

var CommandDataMutex sync.Mutex
var StateValueDataMutex sync.Mutex
var VoteValueDataMutex sync.Mutex
var ConsensusTerminationMutex sync.Mutex
var PQMutex sync.Mutex
var stringIP string
var responceChannelMapMutex sync.Mutex

var responseChannelMap map[string]chan ResponseToClient


var readCnt int = 0
var durationSum time.Duration = 0
var readAverage time.Duration = 0

type Data interface{}

type ConsensusData struct {
	Data Data
}

type CommandData struct {
	Op          string
	Timestamp   time.Time
	Seq         int
	ClientAddr  string
	ReplicaAddr string
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


type ResponseToClient struct {
	Value      int
	ClientAddr string
}

type OpTimestamp struct {
	Op string
	Timestamp time.Time
}

var CommandDataMapList map[int][]CommandData
var StateValueDataMapList map[SeqPhase][]StateValueData
var VoteValueDataMapList map[SeqPhase][]VoteValueData
var ConsensusTerminationMapList map[int][]ConsensusTermination
var Dictionary map[OpTimestamp]bool
var PQ PriorityQueue

func init() {
	CommandDataMapList = make(map[int][]CommandData)
	StateValueDataMapList = make(map[SeqPhase][]StateValueData)
	VoteValueDataMapList = make(map[SeqPhase][]VoteValueData)
	ConsensusTerminationMapList = make(map[int][]ConsensusTermination)
	Dictionary = make(map[OpTimestamp]bool)

	gob.Register(CommandData{})
	gob.Register(StateValueData{})
	gob.Register(VoteValueData{})
	gob.Register(ConsensusTermination{})
	gob.Register(Request{})
	gob.Register(ResponseToClient{})
	gob.Register(ConsensusData{})

	PQ := make(PriorityQueue, 0)
	heap.Init(&PQ)

	selfIP, err := pubip.Get()
	if err != nil {
		fmt.Println("IPアドレスの取得エラー:", err)
		return
	}
	stringIP = selfIP.String()

	responceChannelMapMutex.Lock()
	responseChannelMap = make(map[string]chan ResponseToClient)
	responceChannelMapMutex.Unlock()

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
		return
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {

		//if the connection is closed, return
		_, err := conn.Read(make([]byte, 0))
		if err != nil {
			return
		}

		consensusData, err := receiveData(conn)
		if err != nil {
			return
		}
		data := consensusData.Data

		switch data := data.(type) {
		case CommandData:
			CommandDataMutex.Lock()
			//fmt.Println("Received CommandData: ", data)
			CommandDataMapList[data.Seq] = append(CommandDataMapList[data.Seq], data)
			//fmt.Println("CommandDataMapList: ", CommandDataMapList)
			CommandDataMutex.Unlock()
			return 


		case StateValueData:
			StateValueDataMutex.Lock()
			//fmt.Println("Received StateValueData: ", data)
			StateValueDataMapList[SeqPhase{Seq: data.Seq, Phase: data.Phase}] = append(StateValueDataMapList[SeqPhase{Seq: data.Seq, Phase: data.Phase}], data)
			//fmt.Println("StateValueDataMapList: ", StateValueDataMapList)
			StateValueDataMutex.Unlock()
			return 


		case VoteValueData:
			VoteValueDataMutex.Lock()
			//fmt.Println("Received VoteValueData: ", data)
			VoteValueDataMapList[SeqPhase{Seq: data.Seq, Phase: data.Phase}] = append(VoteValueDataMapList[SeqPhase{Seq: data.Seq, Phase: data.Phase}], data)
			//fmt.Println("VoteValueDataMaplist: ", VoteValueDataMaplist)
			VoteValueDataMutex.Unlock()
			return 

		case ConsensusTermination:
			ConsensusTerminationMutex.Lock()
			//fmt.Println("Received ConsensusTermination: ", data)
			ConsensusTerminationMapList[data.Seq] = append(ConsensusTerminationMapList[data.Seq], data)
			ConsensusTerminationMutex.Unlock()
			return 

		case Request:
			//fmt.Println("Received Request: ", data)

			if data.CommandData.Op[0] == 'R' {


				now := time.Now()
				value, err := parseReadCommand(data.CommandData.Op, StateMachine)
				duration := time.Since(now)

				
				if err == "notFound" {
					response := ResponseToClient{Value: -1}
					sendData(conn, response)
				} else {
					response := ResponseToClient{Value: value}
					sendData(conn, response)
				}

				durationSum += duration
				readCnt++
				readAverage = durationSum / time.Duration(readCnt)
				fmt.Println("Average read time: ", readAverage)



			} else if !data.Redirected {
				data.Redirected = true
				data.CommandData.ReplicaAddr = ownIP
				data.CommandData.ClientAddr = conn.RemoteAddr().String()
				PQMutex.Lock()
				heap.Push(&PQ, &data.CommandData)
				PQMutex.Unlock()
				broadCastData(replicaIPs, data)

				responceChannelMapMutex.Lock()
				if responseChannelMap[data.CommandData.ClientAddr] == nil {
					responseChannelMap[data.CommandData.ClientAddr] = make(chan ResponseToClient)
				}
				responceChannelMapMutex.Unlock()
				//Wait for termination
				go func() {
					//fmt.Println(data.CommandData.ClientAddr)

					responceChannelMapMutex.Lock()
					respChan := responseChannelMap[data.CommandData.ClientAddr]
					responceChannelMapMutex.Unlock()
					// ロックを解放した後、チャネルからデータを受信
					response := <-respChan
					//fmt.Println("Response to client: ", response)
					sendData(conn, response)

				}()
			} else {
				PQMutex.Lock()
				//fmt.Println("received redirected request: " + data.CommandData.Op)
				heap.Push(&PQ, &data.CommandData)
				PQMutex.Unlock()
			}

		default:
			fmt.Println("未知のデータ型です:", data)
		}
	}
}


func broadCastData(IPLists []string, data Data) {
	//remove self IP from IPLists
	var IPs []string
	for _, ip := range IPLists {
		if ip != stringIP{
			IPs = append(IPs , ip)
		}
	}
	//fmt.Println("broadcasting to replicas: ", IPs)
	conns := setConnectionWithOtherReplicas(IPs)
	//fmt.Println("broadcast to replicas: ", conns)
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
