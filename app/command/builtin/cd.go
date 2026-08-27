package builtin

import (
	"fmt"
	"os"
)

func Cd(path string) {
	if path == "~" {
		path = os.Getenv("HOME")
	}
	err := os.Chdir(path)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", path)
		return
	}
}
