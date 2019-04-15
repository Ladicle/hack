package contest

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/PuerkitoBio/goquery"
)

const (
	atCoderBaseURL  = "https://atcoder.jp/contests"
	atCoderQuizPath = "tasks"
)

func NewAtCoder(id string) *AtCoder {
	return &AtCoder{
		ContestID: id,
	}
}

type AtCoder struct {
	ContestID string
}

func (a AtCoder) ContestDir() (string, error) {
	var dir string

	u, err := user.Current()
	if err != nil {
		return dir, err
	}

	// FIXME: contests path need to change
	return filepath.Join(u.HomeDir, config.BaseDir(), "atcoder", a.ContestID), nil
}

func (a AtCoder) QuizDir(quizID string) (string, error) {
	dir, err := a.ContestDir()
	if err != nil {
		return dir, err
	}
	return filepath.Join(dir, quizID), nil
}

func (a AtCoder) QuizzesURL() string {
	return fmt.Sprintf("%v/%v/%v", atCoderBaseURL, a.ContestID, atCoderQuizPath)
}

func (a AtCoder) QuizURL(quizID string) string {
	return fmt.Sprintf("%v/%v", a.QuizzesURL(), quizID)
}

func (a AtCoder) SqrapeQuizzes() ([]string, error) {
	doc, err := goquery.NewDocument(a.QuizzesURL())
	if err != nil {
		return nil, err
	}

	var quizzes []string
	var invalidQuiz bool
	doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
		path, ok := s.Find("a").First().Attr("href")
		if !ok {
			invalidQuiz = true
		}
		// path has /contests/<id>/tasks/<quiz_id>
		section := strings.Split(path, "/")
		quizzes = append(quizzes, section[len(section)-1])
	})

	if invalidQuiz {
		return nil, fmt.Errorf("cannot get valid quizzes from %v", a.QuizzesURL())
	}
	return quizzes, nil
}

func (a AtCoder) SqrapeSample(quizID string) ([]*Sample, error) {
	var ss []*Sample

	doc, err := goquery.NewDocument(a.QuizURL(quizID))
	if err != nil {
		return nil, err
	}

	// AtCoder outputs both JP and EN content in one page.
	doc.Find(".lang-en pre").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			// first block is not sample
			return
		}
		if i%2 == 1 {
			ss = append(ss, &Sample{
				ID:    i/2 + 1,
				Input: s.Text(),
			})
		} else {
			ss[i/2-1].Output = s.Text()
		}
	})
	return ss, nil
}

func (a AtCoder) SetupQuizDir() error {
	quizzes, err := a.SqrapeQuizzes()
	if err != nil {
		return err
	}

	cDir, err := a.ContestDir()
	if err != nil {
		return err
	}

	if err := mkQuizDir(cDir, quizzes); err != nil {
		return err
	}

	for _, quiz := range quizzes {
		ss, err := a.SqrapeSample(quiz)
		if err != nil {
			return err
		}

		qDir, err := a.QuizDir(quiz)
		if err != nil {
			return err
		}

		if err := mkSamples(qDir, ss); err != nil {
			return err
		}
	}
	return nil
}
