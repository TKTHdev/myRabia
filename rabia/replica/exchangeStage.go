package main

import (
	"fmt"
	"net"
	"sync"
	
)

var Mutex sync.Mutex	


func exchangeStage(command Command, portNums []int, port int, ln net.Listener) int{
	conns := setConnectionWithOtherReplicas(portNums, port)
    var state int
    wg := sync.WaitGroup{}
	exchangeAfter(conns, command, &wg)
	state = exchangeBefore(command, portNums, ln, &wg)
    return state 
}

func exchangeBefore(command Command, portNums []int, ln net.Listener, wg *sync.WaitGroup) int {
	var receiveCnt int = 0
	var sameCnt int = 0

	result := make(chan int)
	connCh := make(chan net.Conn)

	go func() {
		for {
			conn, err := ln.Accept()
			defer ln.Close()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				return
			}
			connCh <- conn
		}
	}()

	for {
		select {
		case conn := <-connCh:
			wg.Add(1)
			go func(conn net.Conn) {
				defer wg.Done()
				defer conn.Close()

				receivedCommand, err := receiveCommand(conn)
				if err != nil {
					fmt.Println("Error receiving command: ", err.Error())
					
				}
				Mutex.Lock()
				receiveCnt++
				Mutex.Unlock()
				fmt.Printf("Message received.  Cnt:%d\n", receiveCnt)

				if receivedCommand == command {
					Mutex.Lock()
					sameCnt++
					Mutex.Unlock()
				}

				if receiveCnt == len(portNums)/2+1 {
					if sameCnt == len(portNums)/2+1 {
						result <- 1
					} else {
						result <- 0
					}
				}
			}(conn)
		case res := <-result:
			return res
		}
		if receiveCnt == len(portNums)/2+1 {
			break
		}
	}
	return <-result
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