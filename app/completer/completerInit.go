package completer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/command/builtin"
)

// caching last completed key used
var lastKey = ""

func CompleterInit() *readline.PrefixCompleter {
	completer := readline.NewPrefixCompleter(
		readline.PcItem("exit"),
		readline.PcItem("type"),
		readline.PcItem("echo"),
		readline.PcItem("history"),
		readline.PcItem("jobs"),
		readline.PcItem("declare"),
		readline.PcItem("pwd"),
		readline.PcItem("cd"),
		readline.PcItemDynamic(executableCompletion, nil),
	)
	return completer
}

func executableCompletion(prefixes string) []string {
	suggestions := make([]string, 0)
	prefix := strings.Split(prefixes, " ")
	completeMap := builtin.GetCompleteMap()
	if program, exist := completeMap[prefix[0]]; exist {
		command, current, prev := externalCommand(prefixes)
		cmd := exec.Command(program, command, current, prev)
		cmd.Env = append(cmd.Env,
			"COMP_LINE="+prefixes,
			fmt.Sprintf("COMP_POINT=%d", len(prefixes)),
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			panic(err)
		}
		output := strings.Fields(strings.TrimSpace(string(out)))
		for _, suggestion := range output {
			command = strings.TrimSpace(command)
			prev = strings.TrimSpace(prev)
			if command == prev || prev == "" {
				suggestion = command + " " + suggestion
			} else {
				suggestion = command + " " + prev + " " + suggestion
			}
			suggestions = append(suggestions, suggestion)
		}
		return suggestions
	}
	first, last := prefix[0], prefix[len(prefix)-1]
	if path, exist := completeMap[first]; exist {
		cmd := exec.Command(path)
		lastKey = path
		output, _ := cmd.Output()
		outputString := strings.TrimSpace(string(output))
		suggestions = append(suggestions, prefixes+outputString)
		return suggestions
	}
	if last != "" {
		files, err := os.ReadDir("./")
		if err != nil {
			panic(err)
		}
		prefixes = strings.TrimSpace(prefixes)
		if strings.HasSuffix(prefixes, "/") {
			files, _ = os.ReadDir(prefixes)
		}
		for _, file := range files {
			name := file.Name()
			if strings.HasSuffix(prefixes, "/") {
				name = prefixes + file.Name()
			}
			if file.IsDir() {
				suggestions = append(suggestions, name+"/")
			} else {
				suggestions = append(suggestions, name)
			}
		}
	} else {
		files, err := os.ReadDir("./")
		if err != nil {
			panic(err)
		}
		prefixes = strings.TrimSpace(prefixes)
		if strings.HasSuffix(prefixes, "/") {
			files, _ = os.ReadDir(prefixes)
		}
		for _, file := range files {
			name := file.Name()
			if strings.HasSuffix(prefixes, "/") {
				name = prefixes + file.Name()
			}
			if file.IsDir() {
				suggestions = append(suggestions, prefixes+" "+name+"/")
			} else {
				suggestions = append(suggestions, name)
			}
		}
	}
	if len(prefix) <= 1 {
		path := os.Getenv("PATH")
		dirs := filepath.SplitList(path)
		for _, dir := range dirs {
			files, _ := os.ReadDir(dir)
			for _, file := range files {
				name := file.Name()
				if last == "" {
					name = prefixes + name
				}
				if file.IsDir() {
					suggestions = append(suggestions, name)
				} else {
					suggestions = append(suggestions, name)
				}
			}
		}
	}
	return suggestions
}
func externalCommand(prefix string) (cmd, current, prev string) {
	//ex: cmd prev curre<tab>
	part := strings.Fields(prefix)
	if len(part) == 0 {
		return
	}
	cmd = part[0]
	switch {
	case len(part) == 1:
		current = ""
		prev = ""
	case len(part) == 2:
		current = part[1]
		prev = part[0]
	case len(part) > 2:
		current = part[len(part)-1]
		prev = part[len(part)-2]
	}
	return
}
