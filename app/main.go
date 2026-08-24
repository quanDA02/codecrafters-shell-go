package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
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
		newline, length = c.completer.Do(line, pos)
		// if len(arr) <= 2 {
		// 	newline, length = c.completer.Do(line, pos)
		// 	return newline, length
		// } else {
		// 	fmt.Print("\x07")
		// 	return nil, 0
		// }
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

			// prefix := strings.TrimSpace(string(lline))

			for _, suggestion := range newline {
				fmt.Print(string(arr[len(arr)-1]), string(suggestion), " ")
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

func externalCommand(prefix string) (cmd, current, prev string) {
	//ex: cmd prev curre<tab>
	part := strings.Fields(prefix)
	if len(part) == 0 {
		return
	}
	cmd = part[0]
	switch {
	case len(part) == 1:
		current = ""
		prev = ""
	case len(part) == 2:
		current = part[1]
		prev = part[0]
	case len(part) > 2:
		current = part[len(part)-1]
		prev = part[len(part)-2]
	}
	return
}

// caching last completed key used
var lastKey = ""

func executableCompletion(prefixes string) []string {
	suggestions := make([]string, 0)
	prefix := strings.Split(prefixes, " ")
	if program, exist := completeMap[prefix[0]]; exist {
		command, current, prev := externalCommand(prefixes)
		cmd := exec.Command(program, command, current, prev)
		cmd.Env = append(cmd.Env,
			"COMP_LINE="+prefixes,
			fmt.Sprintf("COMP_POINT=%d", len(prefixes)),
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			panic(err)
		}
		output := strings.Fields(strings.TrimSpace(string(out)))

		// suggestion := output
		for _, suggestion := range output {
			command = strings.TrimSpace(command)
			prev = strings.TrimSpace(prev)
			if command == prev || prev == "" {
				suggestion = command + " " + suggestion
			} else {
				suggestion = command + " " + prev + " " + suggestion
			}
			suggestions = append(suggestions, suggestion)
		}
		// fmt.Println("l:", prefixes)
		// fmt.Println("s:", suggestions)
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
	s = strings.TrimSpace(s)
	builtins := []string{
		"type", "exit", "echo", "complete", "jobs",
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
	key := args[1]
	switch flag {
	case "-p":
		path, exist := completeMap[key]
		if exist {
			fmt.Printf("complete -C '%s' %s\n", path, key)
		} else {
			fmt.Printf("complete: %s: no completion specification\n", key)
		}
	case "-C":
		if len(args) > 2 {
			completeMap[args[2]] = args[1]
		}
	case "-r":
		delete(completeMap, key)
	}
}

// background jobs
type Jobs struct {
	id     int
	name   string
	recent int
	status string
}

var jobMap = make(map[int]*Jobs)

func jobs(doneOnly bool) {
	if len(jobMap) < 1 {
		return
	}
	key := make([]int, 0)
	for _, job := range jobMap {
		key = append(key, job.id)
	}
	sort.Ints(key)
	recents := make([]int, 0)
	for _, job := range jobMap {
		recents = append(recents, job.recent)
	}
	sort.Ints(recents)
	for _, id := range key {
		job := jobMap[id]
		mark := " "
		if job.recent == recents[0] {
			mark = "+"
		}
		if len(jobMap) > 1 && job.recent == recents[1] {
			mark = "-"
		}
		if doneOnly {
			if job.status == "Done" {
				fmt.Printf("[%d]%s  %-24s%s\n", job.id, mark, job.status, job.name)
				delete(jobMap, job.id)
			}
		} else {
			fmt.Printf("[%d]%s  %-24s%s\n", job.id, mark, job.status, job.name)
			if job.status == "Done" {
				delete(jobMap, job.id)
			}
		}
	}
}
func execute(name string) {
	var prevReader *os.File = os.Stdin
	var processes []*exec.Cmd
	commands := pipelineExecute(name)
	for i, command := range commands {
		args, _ := shlex.Split(command)
		isBackground := false
		args, stdout, stderr := redirect(args, os.Stdout, os.Stderr)
		// jobs
		if args[len(args)-1] == "&" {
			isBackground = true
			args = args[0 : len(args)-1]
		}
		cmd := exec.Command(args[0], args[1:]...)
		var prevWriter *os.File
		if i < len(commands)-1 {
			r, w, _ := os.Pipe()
			stdout = w
			prevWriter = w
			prevReader = r
		}
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		cmd.Stdin = prevReader

		// check if it is a built in command
		switch args[0] {
		case "exit":
			os.Exit(0)
			return
		case "type":
			typeCommand(args[1])
			return
		case "echo":
			echo(strings.Join(args[1:], " "), stdout)
		case "complete":
			complete(args[1:])
			return
		case "jobs":
			jobs(false)
			return
		}
		if _, err := exec.LookPath(args[0]); err != nil {
			fmt.Printf("%s: command not found\n", args[0])
			return
		}

		err := cmd.Start()
		if err != nil {
			panic(err)
		}
		if prevWriter != nil {
			prevWriter.Close()
		}
		if isBackground {
			jobID := 1
			for {
				_, exist := jobMap[jobID]
				if exist {
					jobID++
					continue
				}
				break
			}
			job := &Jobs{
				id:     jobID,
				name:   name,
				recent: 0,
				status: "Running",
			}
			jobMap[jobID] = job
			for _, job := range jobMap {
				jobMap[job.id].recent += 1
			}
			fmt.Printf("[%d] %d\n", jobID, cmd.Process.Pid)
			go func(jobID int) {
				cmd.Wait()
				jobMap[jobID].status = "Done"
				jobMap[jobID].name = strings.Join(args, " ")
			}(jobID)
		} else {
			processes = append(processes, cmd)
		}
	}
	for _, cmd := range processes {
		cmd.Wait()
	}
}

func pipelineExecute(name string) (commands []string) {
	commands = strings.Split(name, "|")
	return
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
		jobs(true)
	}
}
