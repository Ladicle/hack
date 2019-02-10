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

// Option manages global command options.
type Option struct {
	WorkDir string
}

// Command run command.
type Command struct {
	Name        string
	Short       string
	Description string
	Run         func(args []string, opt Option) error
}

var cmds []Command

// LoadCmd loads all cmds.
func LoadCmd(io io.Writer) {
	addCmd(NewVersionCmd(io))
	addCmd(NewSetCmd(io))
	addCmd(NewInfoCmd(io))
	addCmd(NewOpenCmd(io))
	addCmd(NewJumpCmd(io))
	addCmd(NewSampleCmd(io))
	addCmd(NewTestCmd(io))
	addCmd(NewCopyCmd(io))
}

func addCmd(c Command) {
	cmds = append(cmds, c)
}

// HandleCmd dispatches and execute specified command.
func HandleCmd(name string, args []string, opt Option) error {
	for _, c := range cmds {
		if name == c.Name {
			return c.Run(args, opt)
		}
	}
	return fmt.Errorf("%q is unknown command", name)
}

// PrintUsage write command usage to specified writer
func PrintUsage() {
	for _, c := range cmds {
		space := strings.Repeat(" ", 25-len(c.Short))
		fmt.Printf("  %s%s%s\n", c.Short, space, c.Description)
	}
}
