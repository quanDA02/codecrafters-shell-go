package builtin

import (
	"fmt"
	"os"
)

func Pwd() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(path)
}
