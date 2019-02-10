package cmd

import (
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
		Short:       "set [CONTEST] [ARGS]...",
		Description: fmt.Sprintf("Set contest information (contests: %v)", strings.Join(contest.ListContestName(), ", ")),
		Run:         s.run,
	}
}

type setCmd struct {
	IO io.Writer
}

func init() {
	contest.LoadContest()
}

func (c *setCmd) run(args []string, opt Option) error {
	if len(args) == 0 {
		return fmt.Errorf("no contest names are specified (contests: %v)",
			strings.Join(contest.ListContestName(), ", "))
	}
	ctt, err := contest.GetContest(args[0])
	if err != nil {
		return fmt.Errorf("%v (contests: %v)",
			err, strings.Join(contest.ListContestName(), ", "))
	}

	if err := ctt.Set(OutputDirectory, args); err != nil {
		return fmt.Errorf("failed to set a contest: %v", err)
	}
	if _, err := fmt.Fprintf(c.IO, "Success to create contest directories to %s\n", OutputDirectory); err != nil {
		return err
	}

	// initialize current quiz pointer
	config.C.CurrentQuizz = ""
	return config.WriteConfig(ConfigPath)
}
