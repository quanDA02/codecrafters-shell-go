package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {

		fmt.Print("$ ")
		command, _ := reader.ReadString('\n')
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
