package contest

import (
	"fmt"
	"path/filepath"
	"strings"

	"flag"

	"github.com/Ladicle/hack/pkg/config"
)

var num int

const alphabet = "abcdefghijklmno"

// NewFreeContest create free contest.
func NewFreeContest() Contest {
	return Contest{
		Name: "free",
		Set:  set,
	}
}

func set(output string, arg []string) error {
	flag.IntVar(&num, "n", 4, "The number of quizzes.")
	flag.IntVar(&num, "number", 4, "The number of quizzes.")
	flag.Parse()
	if flag.NArg() != 1 {
		return fmt.Errorf("invalid argument of number")
	}
	if num > 15 {
		return fmt.Errorf("the quiz number of limit is exceeded (max: 15)")
	}

	baseDir := filepath.Join(output, filepath.Join(arg...))
	quizzes := strings.Split(alphabet[:num], "")
	mkdirs(baseDir, quizzes)

	config.C.Contest = config.Contest{
		Name:    arg[0],
		Path:    baseDir,
		Quizzes: quizzes,
	}
	return nil
}
