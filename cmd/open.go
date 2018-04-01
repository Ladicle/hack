package cmd

import (
	"flag"
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

func (c *openCmd) run() error {
	flag.Parse()
	if flag.NArg() != 0 {
		return fmt.Errorf("invalid number of arguments")
	}
	if err := open.Start(config.C.Contest.URL); err != nil {
		return fmt.Errorf("cloud not open %v in automatically", config.C.Contest.URL)
	}
	return nil
}
