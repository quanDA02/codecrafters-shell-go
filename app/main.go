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

func echo(s string) {
	// tokens, _ := shlex.Split(s)

	// t, _ := shlex.Split(s)
	// var result []string
	// for i, cmd := range t {
	// 	if strings.HasPrefix(cmd, ">") || strings.HasPrefix(cmd, "1>") {
	// 		result = t[1:i]
	// 	}
	// }
	// text := strings.Join(result, " ")
	// for i, cmd := range tokens {
	// 	if cmd == ">" || cmd == "1>" {
	// 		filename := tokens[i+1]
	// 		filename = strings.ReplaceAll(filename, "\"", "")
	// 		//create directory
	// 		err := os.MkdirAll(filepath.Dir(filename), 0750)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		//create file
	// 		err = os.WriteFile(filename, []byte(text+"\n"), 0666)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}

	// 		return
	// 	}
	// 	if cmd == "2>" {
	// 		outputErrorFile, err := os.Create(tokens[len(tokens)-1])
	// 		if err != nil {
	// 			fmt.Fprintf(os.Stderr, "Couldn't create file: %v", err)
	// 		}
	// 		defer outputErrorFile.Close()
	// 		os.Stderr = outputErrorFile
	// 		tokens = tokens[:len(tokens)-2]
	// 		fmt.Println(strings.Join(tokens[1:i], " "))
	// 		return
	// 	}
	// }

	fmt.Printf("%s\n", s)
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

func execute(name string) error {
	// cmd := exec.Command(strings.Join(token, " "))
	var stdout *os.File = os.Stdout

	args, _ := shlex.Split(name)
	if len(args) > 2 && (args[len(args)-2] == ">" || args[len(args)-2] == "1>") {
		outputFile, err := os.Create(args[len(args)-1])
		if err != nil {
			return err
		}
		defer outputFile.Close()
		stdout = outputFile
		args = args[:len(args)-2]
	}

	if len(args) > 2 && args[len(args)-2] == "2>" {
		outputErrorFile, err := os.Create(args[len(args)-1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't create file: %v", err)
		}
		defer outputErrorFile.Close()
		os.Stderr = outputErrorFile
		args = args[:len(args)-2]
	}

	switch args[0] {
	case "exit":
		os.Exit(0)
	case "type":
		typeCommand(name[5:])
	case "echo":
		echo(strings.Join(args[1:], " "))
	}

	if _, err := exec.LookPath(args[0]); err != nil {
		fmt.Printf("%s: command not found\n", args[0])
		return err
	}
	cmd := exec.Command(args[0])
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
		execute(command)

		//check if command : exit
		// switch tokens[0] {
		// case "exit":
		// 	os.Exit(0)
		// case "echo":
		// 	echo(command)
		// case "type":
		// 	typeCommand(command[5:])
		// default:
		// 	execute(command, tokens)
		// }

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
