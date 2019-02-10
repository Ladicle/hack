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
		Name:  atCorder.Name,
		Set:   atCorder.set,
		Usage: "atcorder <contest-type: abc,arc,agc> <contest-number>",
	}
}

func (a *atCorder) set(output string, args []string) error {
	if len(args) < 3 { // args includes this sub-command name(atcoder).
		return fmt.Errorf("contest type and number is required argument, but only got %v", args)
	}

	t, stringN := args[1], args[2]
	in, err := strconv.Atoi(stringN)
	if err != nil {
		return fmt.Errorf("%v is an invalid contest number", stringN)
	}
	n := fmt.Sprintf("%03d", in)

	quizID, err := getQuizID(t)
	if err != nil {
		return err
	}
	quizzes := strings.Split(quizID, "")

	baseDir := filepath.Join(output, a.Name, t, n)
	if err := mkdirs(baseDir, quizzes); err != nil {
		return err
	}

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

func getQuizID(contestType string) (string, error) {
	switch contestType {
	case typeABC:
		return "abcd", nil
	case typeARC:
		return "cdef", nil
	case typeAGC:
		return "abcdef", nil
	}
	return "", fmt.Errorf("%v is unknown AtCorder contest", contestType)
}
