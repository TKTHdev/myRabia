package main

import (
	"fmt"
	"net"
	"sync"
	"time"
	
)

var Mutex sync.Mutex	


func exchangeStage(command Command, portNums []int, port int, ln net.Listener) int{
	conns := setConnectionWithOtherReplicas(portNums, port)
    var state int
    wg := sync.WaitGroup{}
	exchangeAfter(conns, command, &wg)
	state = exchangeBefore(command, portNums, port,ln)
    return state 
}


func exchangeBefore(command Command, portNums []int, port int,ln net.Listener) int {
    var receiveCnt int = 0
    var sameCnt int = 0

    result := make(chan int)

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("接続エラー:", err)
            return -1
        }

        go func(conn net.Conn) {
            defer conn.Close()

            receivedCommand, err := receiveCommand(conn)
            if err != nil {
                return
            }
            receiveCnt++
			fmt.Printf("Message received.  Cnt:%d\n",receiveCnt)
            

            if receivedCommand == command {
                sameCnt++
            }

            if receiveCnt == len(portNums)/2+1 {
                if sameCnt == len(portNums)/2+1 {
                    result <- 1
                } else {
                    result <- 0
                }
            }
        }(conn)


        select {
        case res := <-result:
            return res
        case <-time.After(1 * time.Second):
        }
    }
}

func exchangeAfter(conns []net.Conn, command Command, wg *sync.WaitGroup) {
    for _, conn := range conns {
        wg.Add(1)
        go func(conn net.Conn) {
            defer wg.Done()
            fmt.Printf("Sending command to %v\n", conn.RemoteAddr())
            sendCommand(conn, command)
        }(conn)
    }
}