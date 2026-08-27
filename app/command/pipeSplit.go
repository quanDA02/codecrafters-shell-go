package command

import "strings"

func pipelineSplit(name string) (commands []string) {
	commands = strings.Split(name, "|")
	return
}
