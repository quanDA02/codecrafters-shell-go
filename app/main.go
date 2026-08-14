package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/google/shlex"
)

func echo(s string) {
	tokens := strings.Split(s, " ")

	t, _ := shlex.Split(s)
	var result []string
	for i, cmd := range t {
		if strings.HasPrefix(cmd, ">") {
			result = t[1:i]
		}
	}
	text := strings.Join(result, "")
	for i, cmd := range tokens {
		if cmd == ">" || cmd == "1>" {
			filename := tokens[i+1]
			filename = strings.ReplaceAll(filename, "\"", "")
			//create directory
			err := os.MkdirAll(filepath.Dir(filename), 0750)
			if err != nil {
				log.Fatal(err)
			}
			//create file
			err = os.WriteFile(filename, []byte(text), 0666)
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
}

func findPath(file string) string {
	path, err := exec.LookPath(file)
	if err != nil {
		return ""
	}
	return path
}

func execute(name string, token []string) error {
	// cmd := exec.Command(strings.Join(token, " "))
	var stdout *os.File = os.Stdout

	args := strings.Split(name, " ")
	if len(args) > 2 && (args[len(args)-2] == ">" || args[len(args)-2] == "1>") {
		outputFile, err := os.Create(args[len(args)-1])
		if err != nil {
			panic(err)
		}
		defer outputFile.Close()
		stdout = outputFile
		args = args[:len(args)-2]
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		//literally the strings split but on crack
		tokens, _ := shlex.Split(command)

		//check if command : exit
		switch tokens[0] {
		case "exit":
			os.Exit(0)
		case "echo":
			echo(command)
		case "type":
			typeCommand(command[5:])
		default:
			execute(command, tokens)
		}

		// if command == "exit" {
		// 	break
		// } else if strings.HasPrefix(command, "echo ") {
		// 	echo(command)
		// } else if strings.HasPrefix(command, "type ") {
		// 	typeCommand(command[5:])
		// } else {
		// 	if err := execute(command, tokens); err != nil {
		// 		fmt.Fprintln(os.Stderr, err)
		// 	} else {
		// 		fmt.Println(command + ": command not found")
		// 	}

		// }
	}
}
