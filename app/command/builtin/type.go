package builtin

import (
	"fmt"
	"os/exec"
	"slices"
	"strings"
)

func Type(s string) {
	s = strings.TrimSpace(s)
	builtins := []string{
		"type", "exit", "echo", "complete", "jobs", "history", "declare", "pwd", "cd",
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
