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
		Run: func(args []string, opt Option) error {
			flag.Parse()
			if err := s.validate(args); err != nil {
				return err
			}
			return s.run(args, opt)
		},
	}
}

type jumpCmd struct {
	IO io.Writer
}

func (c *jumpCmd) validate(args []string) error {
	if flag.NArg() >= 2 {
		return fmt.Errorf("invalid number of arguments")
	}
	if len(config.C.Contest.Quizzes) < 1 {
		return fmt.Errorf("You need set contest before")
	}
	if flag.NArg() == 1 && !hasQuiz(flag.Arg(0), config.C.Contest.Quizzes) {
		return fmt.Errorf("%q is unknown quiz", flag.Arg(0))
	}
	return nil
}

func (c *jumpCmd) run(args []string, opt Option) error {
	var quiz string
	switch flag.NArg() {
	case 0:
		quiz = config.C.CurrentQuizz
		if opt.WorkDir == filepath.Join(config.C.Contest.Path, quiz) {
			quiz = nextQuiz(quiz, config.C.Contest.Quizzes)
		}
	case 1:
		quiz = flag.Arg(0)
	}

	if quiz == "" {
		quiz = config.C.Contest.Quizzes[0]
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
		return current
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
