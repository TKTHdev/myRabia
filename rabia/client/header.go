package main

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

func init() {
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
	Redirected  bool
	Timestamp   int
}

func sendData(conn net.Conn, data Data) {
	encoder := gob.NewEncoder(conn)
	err := encoder.Encode(ConsensusData{Data: data})
	if err != nil {
		fmt.Println("データ送信エラー:", err)
		return
	}
}

func generateRandomCommand() string {
	variables := []string{"x", "y", "z", "a", "b", "c"}
	operators := []string{"=", "+", "-", "*", "%"}

	variable := variables[rand.Intn(len(variables))]
	operator := operators[rand.Intn(len(operators))]
	value := rand.Intn(100)

	return fmt.Sprintf("%s %s %d", variable, operator, value)
}

func parseCommand(command string, stateMachine map[string]int) error {
	parts := strings.Split(command, " ")

	if len(parts) == 3 {
		variable := parts[0]
		operator := parts[1]
		value, err := strconv.Atoi(parts[2])
		if err != nil {
			return err
		}

		switch operator {
		case "=":
			stateMachine[variable] = value
		case "+":
			stateMachine[variable] += value
		case "-":
			stateMachine[variable] -= value
		case "*":
			stateMachine[variable] *= value
		case "%":
			if value == 0 {
				value = 1
			}
			stateMachine[variable] %= value
		default:
			return fmt.Errorf("invalid operator: %s", operator)
		}
	} else {
		return fmt.Errorf("invalid command format: %s", command)
	}

	return nil
}
