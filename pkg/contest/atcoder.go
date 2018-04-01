package contest

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Ladicle/hack/pkg/config"
)

const (
	typeABC = "abc" // AtCorder Beginner Contest
	typeARC = "arc" // AtCorder Regular Contest
	typeAGC = "agc" // AtCorder Grand Contest
)

type atCorder struct {
	Name string
}

// NewAtCorderContest create atcoder contest.
func NewAtCorderContest() Contest {
	atCorder := atCorder{Name: "atcoder"}
	return Contest{
		Name: atCorder.Name,
		Set:  atCorder.set,
	}
}

func (a *atCorder) set(output string, arg []string) error {
	if len(arg) != 2 {
		return fmt.Errorf("you need specify AtCorder contest <type> and <number>")
	}

	t := arg[0]
	in, err := strconv.Atoi(arg[1])
	if err != nil {
		return fmt.Errorf("%v is invalid <number> argument", arg[1])
	}
	n := fmt.Sprintf("%03d", in)

	var quizzes []string
	baseDir := filepath.Join(output, a.Name, t, n)
	switch t {
	case typeABC:
		quizzes = strings.Split("abcd", "")
	case typeARC:
		quizzes = strings.Split("cdef", "")
	case typeAGC:
		quizzes = strings.Split("abcdef", "")
	default:
		return fmt.Errorf("%v is unknown AtCorder contest", t)
	}
	mkdirs(baseDir, quizzes)

	config.C.CurrentQuizz = quizzes[0]
	config.C.Contest = config.Contest{
		Name:    a.Name,
		URL:     generateAtCorderURL(t, n),
		Path:    baseDir,
		Quizzes: quizzes,
	}
	return nil
}

func generateAtCorderURL(ctype, number string) string {
	return fmt.Sprintf("https://%s%s.contest.atcoder.jp/", ctype, number)
}
