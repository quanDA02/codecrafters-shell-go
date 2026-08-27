package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/command/builtin"
	"github.com/google/shlex"
)

func Execute(input string) {
	if input == "" {
		return
	}
	builtin.AppendHistory(input)
	var prevReader *os.File = os.Stdin
	var processes []*exec.Cmd
	input = builtin.VariableExpand(input)
	commands := pipelineSplit(input)
	for i, command := range commands {
		command = strings.TrimSpace(command)
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
		var Reader *os.File = os.Stdin
		if i < len(commands)-1 {
			r, w, _ := os.Pipe()
			stdout = w
			prevWriter = w
			Reader = r
		}
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		cmd.Stdin = prevReader
		// check if it is a built in command
		isBuiltin := executeBuiltins(args, stdout)
		if !isBuiltin {
			if _, err := exec.LookPath(args[0]); err != nil {
				fmt.Printf("%s: command not found\n", args[0])
				return
			}
			err := cmd.Start()
			if err != nil {
				panic(err)
			}
		}
		if prevWriter != nil {
			prevWriter.Close()
		}
		if prevReader != os.Stdin {
			prevReader.Close()
		}
		if isBackground {
			builtin.CreateJob(args, cmd, input)
		} else {
			processes = append(processes, cmd)
		}
		prevReader = Reader
	}
	for _, cmd := range processes {
		cmd.Wait()
	}
}
