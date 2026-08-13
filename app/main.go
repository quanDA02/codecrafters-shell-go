package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/google/shlex"
)

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
			echo(command)
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
	tokens := strings.Split(s, " ")
	text := strings.Join(s, " ")
	for i, cmd := range tokens {
		fmt.Println(cmd)
		if cmd == ">" || cmd == "1>" {

			err := os.WriteFile(tokens[i], []byte(text), 0666)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}

	fmt.Println(text)
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
