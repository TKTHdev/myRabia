package main

import (
    "encoding/gob"
    "fmt"
    "net"
    "sync"
)

var CommandDataMutex sync.Mutex
var StateValueDataMutex sync.Mutex
var VoteValueDataMutex sync.Mutex

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
    Value int
    Seq   int
    Phase int
}

type VoteValueData struct {
    Value int
    Seq   int
    Phase int
}

var CommandDataMapList map[int][]CommandData
var StateValueDataMapList map[int][]StateValueData
var VoteValueDataMapList map[int][]VoteValueData

func init() {
    CommandDataMapList = make(map[int][]CommandData)
    StateValueDataMapList = make(map[int][]StateValueData)
    VoteValueDataMapList = make(map[int][]VoteValueData)
    gob.Register(CommandData{})
    gob.Register(StateValueData{})
    gob.Register(VoteValueData{})
}

func listenAndAccept(port string) {
    ln, err := net.Listen("tcp", ":"+port)
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
        //fmt.Println("データ送信エラー:", err)
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
            CommandDataMapList[data.Seq] = append(CommandDataMapList[data.Seq], data)
            //fmt.Println("CommandDataMapList: ", CommandDataMapList)
            CommandDataMutex.Unlock()
        case StateValueData:
            StateValueDataMutex.Lock()
            StateValueDataMapList[data.Seq] = append(StateValueDataMapList[data.Seq], data)
            //fmt.Println("StateValueDataMapList: ", StateValueDataMapList)
            StateValueDataMutex.Unlock()
        case VoteValueData:
            VoteValueDataMutex.Lock()
            VoteValueDataMapList[data.Seq] = append(VoteValueDataMapList[data.Seq], data)
            //fmt.Println("VoteValueDataMaplist: ", VoteValueDataMaplist)
            VoteValueDataMutex.Unlock()
        default:
            fmt.Println("未知のデータ型です:", data)
        }
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

func deleteData(seq int) {
    CommandDataMutex.Lock()
    delete(CommandDataMapList, seq)
    CommandDataMutex.Unlock()

    StateValueDataMutex.Lock()
    delete(StateValueDataMapList, seq)
    StateValueDataMutex.Unlock()

    VoteValueDataMutex.Lock()
    delete(VoteValueDataMapList, seq)
    VoteValueDataMutex.Unlock()
}

//Unused now 
//Measure the size of the map data structure
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
