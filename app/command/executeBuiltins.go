package command

import (
	"os"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/command/builtin"
)

func executeBuiltins(args []string, stdout *os.File) (isBuiltin bool) {
	isBuiltin = true
	switch args[0] {
	case "exit":
		builtin.Exit()
	case "type":
		builtin.Type(args[1])
	case "echo":
		builtin.Echo(strings.Join(args[1:], " "), stdout)
	case "complete":
		builtin.Complete(args[1:])
	case "jobs":
		builtin.Jobs(false)
	case "history":
		builtin.History(args)
	case "declare":
		command := strings.Join(args[1:], " ")
		builtin.Declare(command)
	case "pwd":
		builtin.Pwd()
	case "cd":
		builtin.Cd(strings.Join(args[1:], ""))
	default:
		isBuiltin = false
	}
	return
}
