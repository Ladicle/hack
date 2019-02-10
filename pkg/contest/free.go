package contest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"flag"

	"github.com/Ladicle/hack/pkg/config"
)

var num int

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// NewFreeContest create free contest.
func NewFreeContest() Contest {
	return Contest{
		Name:  "free",
		Set:   set,
		Usage: "free [--number=<number: default(4)>] <path>...",
	}
}

func set(output string, args []string) error {
	parse(args)
	if err := validation(); err != nil {
		return err
	}

	baseDir := filepath.Join(output, filepath.Join(flag.Args()...))
	quizzes := strings.Split(alphabet[:num], "")
	if err := mkdirs(baseDir, quizzes); err != nil {
		return err
	}

	config.C.Contest = config.Contest{
		Name:    "free",
		Path:    baseDir,
		Quizzes: quizzes,
	}
	return nil
}

func parse(args []string) {
	os.Args = args
	flag.IntVar(&num, "n", 4, "The number of quizzes.")
	flag.IntVar(&num, "number", 4, "The number of quizzes.")
	flag.Parse()
}

func validation() error {
	if flag.NArg() < 1 {
		return fmt.Errorf("free command requires at least 1 arguments, but got only %v arguments", flag.NArg())
	}
	if num > 25 {
		return fmt.Errorf("the quiz number of limit is exceeded (max: 25)")
	}
	return nil
}
