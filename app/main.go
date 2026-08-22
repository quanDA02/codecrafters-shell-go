package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	arr, _ := shlex.Split(string(line))
	lline := line
	var newline [][]rune
	var length int
	if _, exist := completeMap[arr[0]]; exist {
		if len(arr) != 2 {
			newline, length = c.completer.Do(line, pos)
			return newline, length
		} else {
			fmt.Print("\x07")
			return nil, 0
		}
	} else {
		//seperate line and only take the last part of it
		if len(arr) > 1 {
			for i, r := range line {
				if r == ' ' {
					lline = line[i:]
					pos = len(lline)
				}
			}
		}
		newline, length = c.completer.Do(lline, pos)
	}
	//these print line for debug purpose only
	// fmt.Println("line:", string(line))
	// fmt.Println("new:", len(newline))
	//bell sound if autocomplete fail
	if len(newline) < 1 {
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

			prefix := strings.TrimSpace(string(lline))

			for _, suggestion := range newline {
				fmt.Print(string(prefix), string(suggestion))
			}
			fmt.Print("\n$ " + string(line))
		}
		return nil, 0
	}

	// checking if there are a slash at the end and remove space trace
	s := strings.TrimSpace(string(newline[0]))
	if strings.HasSuffix(s, "/") {
		newline[0] = newline[0][:len(newline[0])-1]
	}

	return newline, length
}

// check common prefix
func commonPrefix(line [][]rune) []rune {
	first, last := (line[0]), line[len(line)-1]
	result := first[:0]
	for i := 0; i < len(first) && i < len(last) && first[i] == last[i]; i++ {
		result = first[:i+1]
	}
	return result
}

// caching last completed key used
var lastKey = ""

func executableCompletion(prefixes string) []string {
	suggestions := make([]string, 0)
	prefix := strings.Split(prefixes, " ")
	if len(prefix) > 2 {
		if program, exist := completeMap[prefix[0]]; exist {
			command := prefix[0]
			previousWord := prefix[1]
			partial := prefix[2]
			cmd := exec.Command(program, command, partial, previousWord)
			cmd.Env = append(cmd.Env,
				"COMP_LINE="+prefixes,
				fmt.Sprintf("COMP_POINT=%d", len(prefixes)),
			)
			out, err := cmd.Output()
			if err != nil {
				panic(err)
			}
			output := strings.TrimSpace(string(out))
			suggestions = append(suggestions, command+" "+previousWord+" "+output)
		}

		return suggestions
	}
	first, last := prefix[0], prefix[len(prefix)-1]

	// fmt.Println("prefix:", prefix)
	// fmt.Println("fi:", first, "ls:", last)
	// number of line almost double because of some space tracing

	if path, exist := completeMap[first]; exist {
		cmd := exec.Command(path)
		lastKey = path
		output, _ := cmd.Output()
		outputString := strings.TrimSpace(string(output))
		suggestions = append(suggestions, prefixes+outputString)
		return suggestions
	}

	if last != "" {
		files, err := os.ReadDir("./")
		if err != nil {
			panic(err)
		}
		prefixes = strings.TrimSpace(prefixes)
		if strings.HasSuffix(prefixes, "/") {
			files, _ = os.ReadDir(prefixes)
		}
		for _, file := range files {
			name := file.Name()
			if strings.HasSuffix(prefixes, "/") {
				name = prefixes + file.Name()
			}
			if file.IsDir() {
				suggestions = append(suggestions, name+"/")
			} else {
				suggestions = append(suggestions, name)
			}
		}
	} else {
		files, err := os.ReadDir("./")
		if err != nil {
			panic(err)
		}
		prefixes = strings.TrimSpace(prefixes)
		if strings.HasSuffix(prefixes, "/") {
			files, _ = os.ReadDir(prefixes)
		}
		for _, file := range files {
			name := file.Name()
			if strings.HasSuffix(prefixes, "/") {
				name = prefixes + file.Name()
			}
			if file.IsDir() {
				suggestions = append(suggestions, prefixes+" "+name+"/")
			} else {
				suggestions = append(suggestions, name)
			}
		}
	}

	if len(prefix) <= 1 {
		path := os.Getenv("PATH")
		dirs := filepath.SplitList(path)
		for _, dir := range dirs {
			files, _ := os.ReadDir(dir)
			for _, file := range files {
				name := file.Name()
				if last == "" {
					name = prefixes + name
				}
				if file.IsDir() {
					suggestions = append(suggestions, name)
				} else {
					suggestions = append(suggestions, name)
				}
			}
		}
		// fmt.Println(suggestions)
	}

	return suggestions
}

func echo(s string, output *os.File) {
	fmt.Fprintf(output, "%s\n", s)
}

func typeCommand(s string) {
	builtins := []string{
		"type", "exit", "echo", "complete",
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

var completeMap = make(map[string]string)

func complete(args []string) {
	if len(args) < 2 {
		return
	}
	flag := args[0]
	switch flag {
	case "-p":
		key := args[1]
		path, exist := completeMap[key]
		if exist {
			fmt.Printf("complete -C '%s' %s\n", path, key)
		} else {
			fmt.Printf("complete: %s: no completion specification\n", args[1])
		}
	case "-C":
		if len(args) > 2 {
			completeMap[args[2]] = args[1]
		}
	}
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
	case "complete":
		complete(args[1:])
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
		readline.PcItem("complete"),
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
