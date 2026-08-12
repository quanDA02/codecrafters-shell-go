package main

import (
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	// TODO: Uncomment the code below to pass the first stage
	fmt.Print("$ ")
	command := os.Args[1]
	if command != "" {
		fmt.Println(command, ": command not found")
	}
}
