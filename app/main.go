package main

import (
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/command"
	"github.com/codecrafters-io/shell-starter-go/app/command/builtin"
	"github.com/codecrafters-io/shell-starter-go/app/readline"
)

func main() {

	l := readline.ReadlineInit()

	if builtin.HISTFILE != "" {
		builtin.HistoryInit()
	}

	for {
		input, _ := l.Readline()
		input = strings.TrimSpace(input)
		command.Execute(input)
		builtin.Jobs(true)
	}
}
