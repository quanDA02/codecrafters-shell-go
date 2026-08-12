package main

import (
	"fmt"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	// TODO: Uncomment the code below to pass the first stage
	for {
		fmt.Print("$ ")
		var command string
		fmt.Scan(&command)

		command = strings.TrimSpace(command)
		//check if command : exit
		if command == "exit" {
			break
		} else if strings.HasPrefix(command, "echo ") {
			fmt.Println(command[5:])
		}
		// else {
		// 	fmt.Println(command + ": command not found")
		// }
	}

}
