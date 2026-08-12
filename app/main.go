package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
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
		} else if strings.HasPrefix(command, "type ") {
			type_command(command[5:])
		} else {
			fmt.Println(command + ": command not found")
		}
	}
}

func type_command(s string) {
	builtins := []string{
		"type", "exit", "echo",
	}
	if slices.Contains(builtins, s) {
		fmt.Println(s, "is a shell builtin")
	} else {
		path := find_path(s)
		if path != "" {
			fmt.Println(s, "is", path)
		} else {
			fmt.Println(s, "not found")
		}
	}
	// found := false

	// for _, command := range builtins {
	// 	if strings.HasPrefix(s, command) {
	// 		found = true
	// 		fmt.Println(command, "is a shell builtin")
	// 		break
	// 	}
	// }
	// if !found {

	// }
}

func find_path(file string) string {
	path, err := exec.LookPath("string")
	if err != nil {
		return ""
	}
	return path
}
