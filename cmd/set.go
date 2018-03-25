package cmd

import (
	"io"
)

// NewSetCmd sets contest information.
func NewSetCmd(io io.Writer) Command {
	s := setCmd{IO: io}
	return Command{
		Name:        "set",
		Short:       "set PATH",
		Description: "Set contest information",
		Run:         s.run,
	}
}

type setCmd struct {
	IO io.Writer
}

func (c *setCmd) run() error {
	return nil
}
