package runner

import (
	"fmt"
	"strings"
)

type errUnknownCommandExecution struct {
	Err error

	Command string
	Output  string
	Args    []string
}

func (e errUnknownCommandExecution) Error() string {
	var builder strings.Builder

	builder.WriteString("command execution unknown error:\n")
	builder.WriteString(fmt.Sprintf("\t command: '%s'\n", e.Command))
	builder.WriteString(fmt.Sprintf("\t args: '%s'\n", strings.Join(e.Args, " ")))
	builder.WriteString(fmt.Sprintf("\t output: '%s'\n", e.Output))

	return builder.String()
}

func (e errUnknownCommandExecution) Unwrap() error {
	return e.Err
}
