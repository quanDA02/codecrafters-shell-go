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

		tokens := strings.Split(command, " ")

		//check if command : exit
		if command == "exit" {
			break
		} else if strings.HasPrefix(command, "echo ") {
			text := singeQuotes(command[5:])
			fmt.Println(text)
		} else if strings.HasPrefix(command, "type ") {
			typeCommand(command[5:])
		} else {
			path := findPath(tokens[0])
			if path != "" {
				execute(tokens[0], tokens[1:])
			} else {
				fmt.Println(command + ": command not found")
			}

		}
	}
}

func typeCommand(s string) {
	builtins := []string{
		"type", "exit", "echo",
	}
	if slices.Contains(builtins, s) {
		fmt.Println(s, "is a shell builtin")
	} else {
		path := findPath(s)
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

func findPath(file string) string {
	path, err := exec.LookPath(file)
	if err != nil {
		return ""
	}
	return path
}

func execute(name string, args []string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func singeQuotes(input string) string {
	//check for singe quote
	quote := strings.Trim(input, "'")
	return quote
}
