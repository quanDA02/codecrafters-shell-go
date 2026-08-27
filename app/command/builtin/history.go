package builtin

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/shlex"
)

var HISTFILE = os.Getenv("HISTFILE")

// cache history
var cmdHistory = make([]string, 0)
var appendIndex = 0

func HistoryInit() {
	command := "-r " + HISTFILE
	getHistory(command)
	appendIndex = len(cmdHistory)
}

func AppendHistory(input string) {
	cmdHistory = append(cmdHistory, input)
}

func History(args []string) {
	recent := ""
	if len(args) > 1 {
		recent = strings.Join(args[1:], " ")
	}
	getHistory(recent)
}

func getHistory(recent string) {
	index := 0
	args, _ := shlex.Split(recent)
	if len(args) > 0 {
		if _, err := strconv.Atoi(args[0]); err == nil {
			index, _ = strconv.Atoi(args[0])
		}
		switch args[0] {
		case "-r":
			data, err := os.ReadFile(args[1])
			if err != nil {
				panic(err)
			}
			commands := strings.Split((strings.TrimSpace(string(data))), "\n")
			if commands[0] != "" {
				for _, command := range commands {
					cmdHistory = append(cmdHistory, command)
				}
			}
			return
		case "-w":
			file, err := os.Create(args[1])
			if err != nil {
				panic(err)
			}
			defer file.Close()
			text := strings.Join(cmdHistory, "\n") + "\n"
			file.Write([]byte(text))
			return
		case "-a":
			file, err := os.OpenFile(args[1],
				os.O_APPEND|os.O_CREATE|os.O_WRONLY,
				0644)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			if appendIndex < len(cmdHistory) {
				text := strings.Join(cmdHistory[appendIndex:], "\n") + "\n"
				file.Write([]byte(text))
				appendIndex = len(cmdHistory)
			}
			return
		}

	}
	if index == 0 {
		for i, command := range cmdHistory {
			fmt.Printf("%5d %s\n", i+1, command)
		}
	} else {
		for i := len(cmdHistory) - index; i < len(cmdHistory); i++ {
			fmt.Printf("%5d  %s\n", i, cmdHistory[i])
		}
	}
}
