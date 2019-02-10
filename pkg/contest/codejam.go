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
		Name:  codeJam.Name,
		Set:   codeJam.set,
		Usage: "codejam <year> <round>",
	}
}

func (a *codeJam) set(output string, args []string) error {
	if len(args) != 3 { // args contains this sub-command(codejam).
		return fmt.Errorf("contest <year> and <round> is required argument, but only got %v", args)
	}

	year, round := args[1], args[2]
	quizzes := strings.Split("abcd", "")
	baseDir := filepath.Join(output, a.Name, year, round)
	if err := mkdirs(baseDir, quizzes); err != nil {
		return err
	}

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
