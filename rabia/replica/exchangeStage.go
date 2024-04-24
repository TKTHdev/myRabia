package main


import(
	"net"
	"strconv"
	"fmt"
	"sync"
	"time"
)


func exchangeStage(command Command, portNums []int, listener net.Listener, port int)int{
    // 他のレプリカとの同期処理を実装
    var state int
    wg := sync.WaitGroup{}
    wg.Add(1)
    go func() {
        state = exchangeBefore(command, portNums, port)
        wg.Done()
    }()
    go exchangeAfter(portNums, command)
    wg.Wait()
	return state
}

func exchangeBefore(command Command, portNums []int, port int) int {
    ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
    if err != nil {
        fmt.Println("リッスンエラー:", err)
        return -1
    }
    defer ln.Close()

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

func exchangeAfter(portNums []int, command Command) {
	
    for _, portNum := range portNums {
        conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(portNum))
        if err != nil {
			continue	
        }
		fmt.Printf("Sending command to %d\n", portNum)
        defer conn.Close()
        sendCommand(conn, command)
    }
}

