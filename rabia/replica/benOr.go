package main 


import(
	"net"
	"fmt"
	"sync"
	"time"

)


var BenOrMutex sync.Mutex

func roundOne(state int, portNums []int, port int ,ln net.Listener) int {
	conns := setConnectionWithOtherReplicas(portNums, port)
	var vote int
	wg := sync.WaitGroup{}
	vote = roundOneAfter(conns, state)
	wg.Add(1)
	go func() {
		vote = roundOneBefore(portNums,ln)
		wg.Done()
	}()
	wg.Wait()
	return vote
}


func roundOneBefore(portNums []int ,ln net.Listener) int {
	var receiveCnt int = 0
	var oneStateCnt int = 0
	var zeroStateCnt int = 0
	result := make(chan int)
	for{

		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			return -1
		}
		go func(conn net.Conn) {
			receivedState, err := receiveBinaryValue(conn)
			if err != nil {
				return
			}
			BenOrMutex.Lock()
			receiveCnt++
			BenOrMutex.Unlock()
			fmt.Printf("Message received.  Cnt:%d\n", receiveCnt)

			if receivedState == 1 {
				BenOrMutex.Lock()
				oneStateCnt++
				BenOrMutex.Unlock()
			} else {
				BenOrMutex.Lock()
				zeroStateCnt++
				BenOrMutex.Unlock()
			}

			if receiveCnt == len(portNums)/2+1 {
				if oneStateCnt>=len(portNums)/2+1 {
					result <- 1
				}else if zeroStateCnt>=len(portNums)/2+1 {
					result <- 0
				}else {
					result <- -1
				}
			}
		}(conn)

		select {
		case res := <-result:
			return res
		case <-time.After(1 * time.Millisecond):
		}
	
	}

}


func roundOneAfter(conns []net.Conn, state int) int {
	for _, conn := range conns {
		fmt.Printf("Sending state to %v\n", conn.RemoteAddr())
		sendBinaryValue(conn, state)
	}
	return state
}
