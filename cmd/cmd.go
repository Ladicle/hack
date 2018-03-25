package cmd

import (
	"fmt"
	"io"
	"strings"
)

var (
	ConfigPath      string
	OutputDirectory string
)

// Command run command.
type Command struct {
	Name        string
	Short       string
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

// PrintUsage write command usage to specified writer
func PrintUsage(io io.Writer) {
	for _, c := range cmds {
		space := strings.Repeat(" ", 20-len(c.Name))
		fmt.Fprintf(io, "  %s%s%s\n", c.Name, space, c.Description)
	}
}
