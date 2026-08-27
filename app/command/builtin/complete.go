package builtin

import "fmt"

var completeMap = make(map[string]string)

func GetCompleteMap() map[string]string {
	return completeMap
}

func Complete(args []string) {
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
