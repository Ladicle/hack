package cmd

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/contest"
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
	flag.Parse()
	if flag.NArg() >= 2 {
		return fmt.Errorf("invalid number of arguments")
	}

	contest.LoadContest()

	path := strings.Split(flag.Arg(0), "/")
	ctt, err := contest.GetContest(path[0])
	if err != nil {
		return err
	}

	if err := ctt.Set(OutputDirectory, path[1:]); err != nil {
		return err
	}

	return config.WriteConfig(ConfigPath)
}
