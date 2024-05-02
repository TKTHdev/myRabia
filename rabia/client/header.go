package main 

import(
	"net"
	"fmt"
	"encoding/gob"
)


func init(){
	gob.Register(Request{})

}

type Data interface{}

type ConsensusData struct {
    Data Data
}

type CommandData struct {
    Op        string
    Timestamp int
    Seq       int
}



type Request struct {
    CommandData CommandData
    Redirected bool
    Timestamp int
}



func sendData(conn net.Conn, data Data) {
    encoder := gob.NewEncoder(conn)
    err := encoder.Encode(ConsensusData{Data: data})
    if err != nil {
        fmt.Println("データ送信エラー:", err)
        return
    }
}