package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	//"math/rand"

)

var wg sync.WaitGroup
var listener net.Listener
var replicaIPs []string

func main() {
	// Set the number of replicas
	var n int
	fmt.Println("Enter the number of replicas: ")
	_, err := fmt.Scan(&n)
	if err != nil {
		return
	}

	// Initialize the server
	listener = initServer()
	if listener == nil {
		fmt.Println("Server initialization failed")
		return
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("クローズエラー:", err)
		}
	}(listener)

	// Accept connections from n replicas
	acceptNConnections(listener, n)


	fmt.Println("Internet disconnection test? [Y]/[N]")
	var disconnectionTest string
	_, err = fmt.Scan(&disconnectionTest)
	if err != nil {
		return
	}
	if disconnectionTest == "Y" {
		fmt.Println("How many disconnections do you want to test?")
		var num int
		_, err = fmt.Scan(&num)
		if err != nil {
			return
		}
		sendPortNumListToReplicasWithDisconnection(num)		
	}else{
		sendPortNumListToReplicas()
	}
	

}

func initServer() net.Listener {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("リッスンエラー:", err)
		return nil
	}

	fmt.Println("サーバーが起動しました。クライアントからの接続を待機しています...")
	return listener
}

func acceptNConnections(listener net.Listener, n int) {
	var cnt = 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("接続エラー:", err)
			continue
		}
		fmt.Println("接続しました")
		wg.Add(1)
		go handleConnection(conn)
		cnt++
		if cnt == n {
			fmt.Println("All replicas connected")
			break
		}
	}
	wg.Wait()
}

func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("クローズエラー:", err)
		}
	}(conn)
	//get replicas public IP

	fmt.Println("Connected to replica: ", conn.RemoteAddr().String())

	replicaIPs = append(replicaIPs, removePort(conn.RemoteAddr().String()))
	wg.Done()
}

// Send the list of port numbers to the replicas
func sendPortNumListToReplicas() {
	for _, IP := range replicaIPs {
		wg.Add(1)
		go func(IP string) {
			conn, err := net.Dial("tcp", IP+":8080")
			if err != nil {
				fmt.Println("接続エラー:", err)
				return
			}
			defer func(conn net.Conn) {
				err := conn.Close()
				if err != nil {
					fmt.Println("クローズエラー:", err)
				}
			}(conn)
			fmt.Println("Sending port number list to replica")
			_, err = fmt.Fprintf(conn, portListToString(replicaIPs)+"\n")
			if err != nil {
				return
			}
			fmt.Printf("Port number list sent to replica %s\n", IP)
			wg.Done()
		}(IP)
	}
	wg.Wait()
}

func sendPortNumListToReplicasWithDisconnection(num int) {
	//make list of replicas for each replicas
	var replicaIPMap = make(map[string][]string)
	for _, IP := range replicaIPs {
		replicaIPMap[IP] = replicaIPs
	}	
	//choose random replicas to disconnect
	for i := 0; i < num; i++ {
		var IP1, IP2 string
		for{
			IP1 = replicaIPs[0]
			IP2 = replicaIPs[1]
			if IP1 != IP2 {
				break
			}
		}
		//disconnect IP1 and IP2
		replicaIPMap = removeIPs(replicaIPMap, IP1, IP2)
	}

	fmt.Println("replicaIPMap: ", replicaIPMap)

	//send port number list to replicas
	for IP, _ := range replicaIPMap {
		wg.Add(1)
		fmt.Println("IP: ", IP)
		go func(IP string) {
			conn, err := net.Dial("tcp", IP+":8080")
			if err != nil {
				fmt.Println("接続エラー:", err)
				return
			}
			defer func(conn net.Conn) {
				err := conn.Close()
				if err != nil {
					fmt.Println("クローズエラー:", err)
				}
			}(conn)
			fmt.Println("Sending port number list to replica")
			_, err = fmt.Fprintf(conn, portListToString(replicaIPMap[IP])+"\n")
			if err != nil {
				return
			}
			fmt.Printf("Port number list sent to replica %d\n", net.IPv4len)
			wg.Done()
		}(IP)
	}
	wg.Wait()
}


func removeIPs(replicaIPMap map[string][]string, IP1 string , IP2 string) map[string][]string {
	for _, IP := range replicaIPMap[IP1] {
		if IP == IP2 {
			replicaIPMap[IP1] = removeIP(replicaIPMap[IP1], IP)
			break
		}
	}

	for _, IP := range replicaIPMap[IP2] {
		if IP == IP1 {
			replicaIPMap[IP2] = removeIP(replicaIPMap[IP2], IP)
			break
		}
	}
	return replicaIPMap
}


func removeIP(IPs []string, IP string) []string {
	//remove IP from IPs
	var newIPs []string
	for _, ip := range IPs {
		if ip != IP {
			newIPs = append(newIPs, ip)
		}else{
			newIPs = append(newIPs, "N"+ip)
		}
	}
	return newIPs
}

// Convert the list of port numbers to a string
func portListToString(IPs []string) string {
	var IPListString string
	for i, IP := range IPs {
		IPListString += IP
		if i < len(replicaIPs)-1 {
			IPListString += ","
		}
	}
	return IPListString
}

func removePort(address string) string {
	// ":"の位置を見つける
	colonIndex := strings.LastIndex(address, ":")
	if colonIndex == -1 {
		// ":"が見つからなければ、そのまま返す
		return address
	}
	// ":"以前の部分を切り取る
	return address[:colonIndex]
}
