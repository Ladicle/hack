package cmd

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"

	"github.com/Ladicle/hack/pkg/config"
)

// NewJumpCmd jumps contest jumprmation.
func NewJumpCmd(io io.Writer) Command {
	s := jumpCmd{IO: io}
	return Command{
		Name:        "jump",
		Short:       "jump [QUIZ]",
		Description: "Jump move to quiz directory. if not specified quiz option, move to next quiz directory.",
		Run:         s.run,
	}
}

type jumpCmd struct {
	IO io.Writer
}

func (c *jumpCmd) run(args []string, opt Option) error {
	var quiz string

	flag.Parse()
	switch flag.NArg() {
	case 0:
		current := config.C.CurrentQuizz
		if opt.WorkDir != filepath.Join(config.C.Contest.Path, current) {
			quiz = current
		} else {
			quiz = nextQuiz(current, config.C.Contest.Quizzes)
		}
		if quiz == "" {
			return fmt.Errorf("%q is a last quiz, so has not a next quiz", current)
		}
	case 1:
		quiz = flag.Arg(0)
		if !hasQuiz(quiz, config.C.Contest.Quizzes) {
			return fmt.Errorf("%q is unknown quiz", quiz)
		}
	default:
		return fmt.Errorf("invalid number of arguments")
	}

	config.C.CurrentQuizz = quiz
	if err := config.WriteConfig(ConfigPath); err != nil {
		return fmt.Errorf("could not update configuration: %v", err)
	}

	dir := filepath.Join(config.C.Contest.Path, quiz)
	fmt.Fprintf(c.IO, dir)
	return nil
}

func hasQuiz(quiz string, list []string) bool {
	for _, q := range list {
		if q == quiz {
			return true
		}
	}
	return false
}

func nextQuiz(current string, list []string) string {
	if current == "" {
		return list[0]
	}
	var flag bool
	for _, q := range list {
		if flag {
			return q
		}
		if q == current {
			flag = true
		}
	}
	return list[len(list)-1]
}
