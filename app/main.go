package main

import (
	"fmt"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	// TODO: Uncomment the code below to pass the first stage
	for {
		fmt.Print("$ ")
		var command string
		fmt.Scan(&command)

		//check if command : exit
		if command == "exit" {
			break
		}

		command = command + ": command not found"
		fmt.Println(command)
	}

}
