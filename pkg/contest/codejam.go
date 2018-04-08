package contest

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Ladicle/hack/pkg/config"
)

type codeJam struct {
	Name string
}

// NewCodeJamContest create codejam contest.
func NewCodeJamContest() Contest {
	codeJam := codeJam{Name: "codejam"}
	return Contest{
		Name: codeJam.Name,
		Set:  codeJam.set,
	}
}

func (a *codeJam) set(output string, arg []string) error {
	if len(arg) != 2 {
		return fmt.Errorf("you need specify CodeJam contest <year> and <round>")
	}

	year, round := arg[0], arg[1]
	quizzes := strings.Split("abcd", "")
	baseDir := filepath.Join(output, a.Name, year, round)
	mkdirs(baseDir, quizzes)

	config.C.Contest = config.Contest{
		Name:    a.Name,
		URL:     generateCodeJamURL(year),
		Path:    baseDir,
		Quizzes: quizzes,
	}
	return nil
}

func generateCodeJamURL(year string) string {
	return fmt.Sprintf("https://codejam.withgoogle.com/%v", year)
}
