package completer

import (
	"fmt"
	"slices"
	"strings"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/command/builtin"
	"github.com/google/shlex"
)

// wrap autocompele from chzyer/readline and add bell sound
type CompleterBell struct {
	Completer readline.AutoCompleter
	TabCount  int
}

func (c *CompleterBell) Do(line []rune, pos int) ([][]rune, int) {
	arr, _ := shlex.Split(string(line))
	lline := line
	var newline [][]rune
	var length int
	completeMap := builtin.GetCompleteMap()
	if _, exist := completeMap[arr[0]]; exist {
		newline, length = c.Completer.Do(line, pos)
	} else {
		//seperate line and only take the last part of it
		if len(arr) > 1 {
			for i, r := range line {
				if r == ' ' {
					lline = line[i:]
					pos = len(lline)
				}
			}
		}
		newline, length = c.Completer.Do(lline, pos)
	}
	//bell sound if autocomplete fail
	if len(newline) < 1 {
		fmt.Print("\x07")
		return nil, 0
	}
	if len(newline) > 1 {
		slices.SortFunc(newline, slices.Compare)
		b := CommonPrefix(newline)
		if len(b) > 0 {
			return newline, length
		}
		if c.TabCount == 0 {
			fmt.Print("\x07")
			c.TabCount++
		} else {
			c.TabCount = 0
			fmt.Println()
			// print sorted suggestions
			// prefix := strings.TrimSpace(string(lline))
			for _, suggestion := range newline {
				fmt.Print(string(arr[len(arr)-1]), string(suggestion), " ")
			}
			fmt.Print("\n$ " + string(line))
		}
		return nil, 0
	}
	// checking if there are a slash at the end and remove space trace
	s := strings.TrimSpace(string(newline[0]))
	if strings.HasSuffix(s, "/") {
		newline[0] = newline[0][:len(newline[0])-1]
	}
	return newline, length
}
