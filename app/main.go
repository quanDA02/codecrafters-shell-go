package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/google/shlex"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {

		fmt.Print("$ ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		//literally the strings split but on crack
		tokens, _ := shlex.Split(command)
		//check if command : exit
		if command == "exit" {
			break
		} else if strings.HasPrefix(command, "echo ") {
			fmt.Println(strings.Join(tokens[1:], " "))
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

func echo(s string) {

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

func filter(s string) string {
	var tokens []string
	if strings.HasPrefix(s, "'") {
		start := strings.Index(s, "'")
		end := strings.Index(s[start+1:], "'")
		tokens = append(tokens, s[start+1:start+1+end])
	} else if strings.HasPrefix(s, "'") {
	}

	result := strings.Join(tokens, "")
	return result
}

func singeQuotes(s string) string {
	//check for singe quote
	quote := strings.Trim(s, "'")
	return quote
}
