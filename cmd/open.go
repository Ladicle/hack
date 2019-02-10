package cmd

import (
	"fmt"
	"io"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/skratchdot/open-golang/open"
)

// NewOpenCmd opens contest openrmation.
func NewOpenCmd(io io.Writer) Command {
	s := openCmd{IO: io}
	return Command{
		Name:        "open",
		Short:       "open",
		Description: "Open shows contest page in browser",
		Run:         s.run,
	}
}

type openCmd struct {
	IO io.Writer
}

func (c *openCmd) run(args []string, opt Option) error {
	if config.C.Contest.URL == "" {
		return fmt.Errorf("this contest has no URL")
	}
	return open.Start(config.C.Contest.URL)
}
