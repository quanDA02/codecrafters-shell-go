package readline

import (
	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/completer"
)

func ReadlineInit() *readline.Instance {
	autoComplete := completer.CompleterInit()
	l, err := readline.NewEx(&readline.Config{
		Prompt: "$ ",
		AutoComplete: &completer.CompleterBell{
			Completer: autoComplete,
			TabCount:  0},
	})
	if err != nil {
		panic(err)
	}
	return l
}
