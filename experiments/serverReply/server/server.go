package main

import (
    "encoding/gob"
    "fmt"
    "net"
)

type Person struct {
    Name string
    Age  int
}

func main() {
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("Error listening:", err)
        return
    }
    defer listener.Close()

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }

        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()

    decoder := gob.NewDecoder(conn)
    var person Person
    err := decoder.Decode(&person)
    if err != nil {
        fmt.Println("Error decoding:", err)
        return
    }

    fmt.Printf("Received person: %+v\n", person)

    person.Age += 1
    encoder := gob.NewEncoder(conn)
    err = encoder.Encode(person)
    if err != nil {
        fmt.Println("Error encoding:", err)
    }
}