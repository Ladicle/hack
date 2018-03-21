package cmd

import (
	"fmt"
	"io"
)

// Command run command.
type Command struct {
	Name        string
	Description string
	Run         func() error
}

var cmds []Command

// LoadCmd loads all cmds.
func LoadCmd(io io.Writer) {
	addCmd(NewVersionCmd(io))
	addCmd(NewAtCorderCmd(io))
}

func addCmd(c Command) {
	cmds = append(cmds, c)
}

// HandleCmd dispatches and execute specified command.
func HandleCmd(name string) error {
	for _, c := range cmds {
		if name == c.Name {
			return c.Run()
		}
	}
	return fmt.Errorf("%v is unknown command", name)
}
