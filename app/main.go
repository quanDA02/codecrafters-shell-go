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

func echo(s string, output *os.File) {
	fmt.Fprintf(output, "%s\n", s)
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

func execute(name string) {
	stdout := os.Stdout

	args, _ := shlex.Split(name)
	//operators checking
	//redirect stdout to file
	//create
	if len(args) > 2 && (args[len(args)-2] == ">" || args[len(args)-2] == "1>") {
		outputFile, err := os.Create(args[len(args)-1])
		if err != nil {
			return
		}
		defer outputFile.Close()
		stdout = outputFile
		args = args[:len(args)-2]
	}
	//error
	if len(args) > 2 && args[len(args)-2] == "2>" {
		outputErrorFile, err := os.Create(args[len(args)-1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't create file: %v", err)
		}
		defer outputErrorFile.Close()
		os.Stderr = outputErrorFile
		args = args[:len(args)-2]
	}
	//append
	if len(args) > 2 && (args[len(args)-2] == ">>" || args[len(args)-2] == "1>>") {
		outputFile, err := os.OpenFile(args[len(args)-1],
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644)
		if err != nil {
			return
		}
		defer outputFile.Close()
		stdout = outputFile
		args = args[:len(args)-2]
	}
	//append error
	if len(args) > 2 && args[len(args)-2] == "2>>" {
		outputErrorFile, err := os.OpenFile(args[len(args)-1],
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't create file: %v", err)
		}
		defer outputErrorFile.Close()
		os.Stderr = outputErrorFile
		args = args[:len(args)-2]
	}
	// check if it is a built in command
	switch args[0] {
	case "exit":
		os.Exit(0)
		return
	case "type":
		typeCommand(name[5:])
		return
	case "echo":
		echo(strings.Join(args[1:], " "), stdout)
		return
	}

	if _, err := exec.LookPath(args[0]); err != nil {
		fmt.Printf("%s: command not found\n", args[0])
		return
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)
		execute(command)
	}
}
