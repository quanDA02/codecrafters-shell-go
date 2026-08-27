package builtin

import (
	"os"
)

func Exit() {
	if HISTFILE != "" {
		getHistory("-a " + HISTFILE)
	}
	os.Exit(0)
}
