package main

import (
	"fmt"
	"strconv"
	"strings"
)

func parseWriteCommand(command string, stateMachine map[string]int) error {
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

func parseReadCommand(command string, stateMachine map[string]int) (int, string) {
	parts := strings.Split(command, " ")
	fmt.Println(stateMachine)
	if len(parts) == 2 {
		key := parts[1]
		fmt.Println("Key, Value : ", key, stateMachine[key])
		value, ok := stateMachine[key]
		if !ok {
			fmt.Println("Variable not found: ", key)
			return 0, "notFound"
		}
		return value, ""
	} else {
		fmt.Println("Invalid command format: ", command)
		return 0, "invalidFormat"	
	}
}