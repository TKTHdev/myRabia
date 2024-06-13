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
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Error dialing:", err)
        return
    }
    defer conn.Close()

    person := Person{Name: "John", Age: 30}
    encoder := gob.NewEncoder(conn)
    err = encoder.Encode(person)
    if err != nil {
        fmt.Println("Error encoding:", err)
        return
    }

    decoder := gob.NewDecoder(conn)
    var receivedPerson Person
    err = decoder.Decode(&receivedPerson)
    if err != nil {
        fmt.Println("Error decoding:", err)
        return
    }

    fmt.Printf("Received person: %+v\n", receivedPerson)
}