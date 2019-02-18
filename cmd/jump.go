package cmd

import (
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
	if len(config.C.Contest.Quizzes) < 1 {
		return fmt.Errorf("no contest is set")
	}
	return nil
}

func (c *jumpCmd) run(args []string, opt Option) error {
	var quiz string
	if len(args) > 1 {
		// set specified quiz if it exists
		inQuiz := args[1]
		if !hasQuiz(inQuiz, config.C.Contest.Quizzes) {
			return fmt.Errorf("%q is unknown quiz", inQuiz)
		}
		quiz = inQuiz
	} else {
		quiz = config.C.CurrentQuizz
		// set next quiz if the current directory is a quiz directory
		if opt.WorkDir == filepath.Join(config.C.Contest.Path, quiz) {
			quiz = nextQuiz(quiz, config.C.Contest.Quizzes)
		}
		// set the first quiz for the first time
		if quiz == "" {
			quiz = config.C.Contest.Quizzes[0]
		}
	}

	config.C.CurrentQuizz = quiz
	if err := config.WriteConfig(ConfigPath); err != nil {
		return fmt.Errorf("failed to update configuration: %v", err)
	}

	dir := filepath.Join(config.C.Contest.Path, quiz)
	_, err := fmt.Fprintf(c.IO, dir)
	return err
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
