package cmd

import (
	"fmt"
	"io"
)

var (
	// These variables are set at the build time.
	version string
	gitRepo string
)

// NewVersionCmd create version command.
func NewVersionCmd(io io.Writer) Command {
	return Command{
		Name:        "version",
		Short:       "version",
		Description: "Show this command version",
		Run: func(args []string, opt Option) error {
			_, err := fmt.Fprintf(io, "%v version is %v", gitRepo, version)
			return err
		},
	}
}
