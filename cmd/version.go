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
		Description: "Show this command version",
		Run: func() error {
			fmt.Fprintf(io, "%v version is %v", gitRepo, version)
			return nil
		},
	}
}
