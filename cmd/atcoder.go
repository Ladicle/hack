package cmd

import (
	"flag"
	"fmt"
	"io"
)

// NewAtCorderCmd create atcoder command.
func NewAtCorderCmd(io io.Writer) Command {
	a := AtCorder{
		IO: io,
	}
	return Command{
		Name:        "atcoder",
		Description: "Support AtCoder contests",
		Run:         a.run,
	}
}

// AtCorder support AtCorder contests.
type AtCorder struct {
	IO io.Writer
}

func (a *AtCorder) run() error {
	flag.Usage = func() {
		fmt.Fprintln(a.IO, `Usage: hack atcoder [OPTIONS] COMMAND

Options:
  -h --help           Show this help message

Commands:
  add                 Add test case for specified contests
  info                Show current contests information
  jump [QUIZ]         Move to specified quiz directory
  set  TYPE NUMBER    Create source directory and jump to there
  test                Tests program

Types:
  abc                 AtCorder Beginner Contest
  arc                 AtCorder Regular Contest
  atc                 AtCorder Typical Contest`)
	}

	flag.Parse()
	if flag.NArg() == 0 {
		return fmt.Errorf("invalid number of arguments")
	}

	switch flag.Arg(0) {
	case "add":
	case "info":
	case "jump":
	case "test":
	case "set":
	default:
		return fmt.Errorf("%v is unknown command", flag.Arg(0))
	}
	return nil
}
