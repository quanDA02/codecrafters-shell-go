package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/chzyer/readline"
	"github.com/google/shlex"
)

// wrap autocompele from chzyer/readline and add bell sound
type completerBell struct {
	completer readline.AutoCompleter
}

// modified autocomplete Do and add bell sound
func (c *completerBell) Do(line []rune, pos int) ([][]rune, int) {
	newline, length := c.completer.Do(line, pos)
	//bell sound if autocomplete fail
	if len(newline) == 0 {
		fmt.Print("\x07")
	}
	return newline, length
}
func executableCompletion(prefix string) []string {
	path := os.Getenv("PATH")
	dirs := strings.Split(path, ":")
	if path == "" {
		// fmt.Print("fuck")
	}
	fmt.Print(dirs[2] + "ok")
	str := []string{
		"cusssstom_exe_7091", "hello",
	}
	return str
}

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

func redirect(args []string, stdout, stderr *os.File) ([]string, *os.File, *os.File) {
	if len(args) <= 2 {
		return args, stdout, stderr
	}
	operator := args[len(args)-2]
	filename := args[len(args)-1]
	//checking operator
	switch operator {

	// overwrite
	case ">", "1>":
		outputFile, _ := os.Create(filename)
		stdout = outputFile
	case "2>":
		outputFile, _ := os.Create(filename)
		stderr = outputFile
	// append
	case ">>", "1>>":
		outputFile, _ := os.OpenFile(args[len(args)-1],
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644)
		stdout = outputFile
	case "2>>":
		outputFile, _ := os.OpenFile(args[len(args)-1],
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644)
		stderr = outputFile
	default:
		return args, stdout, stderr
	}
	args = args[:len(args)-2]
	return args, stdout, stderr
}

func execute(name string) {

	args, _ := shlex.Split(name)

	args, stdout, stderr := redirect(args, os.Stdout, os.Stderr)

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
	cmd.Stderr = stderr
	cmd.Run()
}

func main() {

	completer := readline.NewPrefixCompleter(
		readline.PcItem("exit"),
		readline.PcItem("type"),
		readline.PcItem("echo"),
		readline.PcItemDynamic(executableCompletion, nil),
	)

	l, err := readline.NewEx(&readline.Config{
		Prompt:       "$ ",
		AutoComplete: &completerBell{completer},
	})
	if err != nil {
		fmt.Println(err)
	}
	for {
		command, _ := l.Readline()
		command = strings.TrimSpace(command)

		execute(command)
	}
}
