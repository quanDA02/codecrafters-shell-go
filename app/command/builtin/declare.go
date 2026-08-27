package builtin

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/shlex"
)

var variables = make(map[string]string)

func VariableExpand(input string) string {
	vars := regexp.MustCompile(`\$(({)?[a-zA-Z_]\w*(})?)`).FindAllString(input, -1)
	for _, cmd := range vars {
		v := strings.TrimPrefix(cmd, "$")
		v = strings.TrimPrefix(v, "{")
		v = strings.TrimSuffix(v, "}")
		value, _ := variables[v]
		input = strings.ReplaceAll(input, cmd, value)
	}
	return input
}

func Declare(command string) {
	args, _ := shlex.Split(command)
	if len(args) < 1 {
		return
	}
	switch args[0] {
	case "-p":
		key := args[1]
		if value, exist := variables[key]; exist {
			fmt.Printf("declare -- %s=\"%s\"\n", key, value)
		} else {
			fmt.Printf("declare: %s: not found\n", args[1])
		}
	default:

		s := strings.Split(strings.TrimSpace(args[0]), "=")
		key, value := s[0], s[1]
		if !regexp.MustCompile(`^[a-zA-Z_]\w*$`).MatchString(key) {
			fmt.Printf("declare: `%s=%s': not a valid identifier\n", key, value)
			return
		}
		variables[key] = value
	}
}
