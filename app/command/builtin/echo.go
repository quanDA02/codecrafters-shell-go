package builtin

import (
	"fmt"
	"os"
)

func Echo(s string, output *os.File) {
	fmt.Fprintf(output, "%s\n", s)
}
