package main 

import(
	"fmt"
	"strings"
	"strconv"
)



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
				value=1		
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
