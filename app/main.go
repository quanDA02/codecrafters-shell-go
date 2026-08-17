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
	tabCount  int
}

// modified autocomplete Do and add bell sound
func (c *completerBell) Do(line []rune, pos int) ([][]rune, int) {
	lline := line
	arr, _ := shlex.Split(string(line))

	//seperate line and only take the last part of it
	if len(arr) == 1 {
		for i, r := range line {
			if r == ' ' {
				lline = line[i:]
				pos = len(lline)
			}
		}
	}
	newline, length := c.completer.Do(lline, pos)

	//these print line for debug purpose only
	// fmt.Println("line:", string(lline))
	// fmt.Println("new:", newline)

	//bell sound if autocomplete fail
	if len(newline) == 0 {
		fmt.Print("\x07")
		return nil, 0
	}

	if len(newline) > 1 {
		slices.SortFunc(newline, slices.Compare)
		b := commonPrefix(newline)
		if len(b) > 0 {
			return newline, length
		}

		if c.tabCount == 0 {
			fmt.Print("\x07")
			c.tabCount++
		} else {
			c.tabCount = 0
			fmt.Println()
			// print sorted suggestions
			for _, suggestion := range newline {
				fmt.Print(string(line), string(suggestion))
			}

			fmt.Print("\n$ " + string(line))
		}
		return nil, 0
	}

	return newline, length
}

// check common prefix
func commonPrefix(line [][]rune) []rune {
	first, last := line[0], line[len(line)-1]
	result := first[:0]
	for i := 0; i < len(first) && i < len(last) && first[i] == last[i]; i++ {
		result = first[:i+1]
	}
	return result
}

func executableCompletion(prefixes string) []string {

	prefix := strings.Split(prefixes, " ")
	first, last := prefix[0], prefix[len(prefix)-1]

	// fmt.Println("prefix:", prefix)
	// fmt.Println("fi:", first, "ls:", last)

	suggestions := make([]string, 0)
	if first == last {

		path := os.Getenv("PATH")
		dirs := strings.Split(path, ":")

		for _, dir := range dirs {
			files, _ := os.ReadDir(dir)
			for _, file := range files {
				suggestions = append(suggestions, ""+file.Name())
			}
		}
	} else {
		files, _ := os.ReadDir("./")
		if strings.HasSuffix(last, "/") {
			files, _ = os.ReadDir(last)
		}

		for _, file := range files {
			name := file.Name()
			//check if the last character is a "/"(slash)
			if strings.HasPrefix(name, last) {
				suggestions = append(suggestions, file.Name())
			} else {
				suggestions = append(suggestions, last+file.Name())
			}
		}
	}

	fmt.Println("suggestion :", suggestions)

	return suggestions
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
		AutoComplete: &completerBell{completer, 0},
	})
	if err != nil {
		panic(err)
	}

	for {
		command, _ := l.Readline()
		command = strings.TrimSpace(command)

		execute(command)
	}
}
