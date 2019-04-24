package contest

import (
	"fmt"
	"net/http"
	"net/url"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/httputil"
	"github.com/PuerkitoBio/goquery"
)

const (
	atCoderBaseURL  = "https://atcoder.jp/contests"
	atCoderLoginURL = "https://atcoder.jp/login"

	atCoderQuizPath = "tasks"
)

func NewAtCoder(id string) *AtCoder {
	return &AtCoder{
		ContestID: id,
		Session:   &httputil.Session{},
	}
}

type AtCoder struct {
	ContestID string
	Session   *httputil.Session
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

func (a AtCoder) getCsrfToken(url string) (string, error) {
	resp, err := a.Session.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", err
	}

	var csrfToken string
	doc.Find("input[name=\"csrf_token\"]").Each(func(i int, s *goquery.Selection) {
		value, exist := s.Attr("value")
		if exist {
			csrfToken = value
		}
	})
	if csrfToken == "" {
		return "", fmt.Errorf("failed to scrape csrfToken from %v", url)
	}
	return csrfToken, nil
}

func (a AtCoder) Login() error {
	// allready logind
	if a.Session.Cookies != nil {
		return nil
	}

	csrfToken, err := a.getCsrfToken(atCoderLoginURL)
	if err != nil {
		return err
	}
	values := url.Values{}
	values.Add("username", config.AtCoderUser())
	values.Add("password", config.AtCoderPass())
	values.Add("csrf_token", csrfToken)

	resp, err := a.Session.PostForm(atCoderLoginURL, &values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (a AtCoder) loginAndGet(url string) (*http.Response, error) {
	if err := a.Login(); err != nil {
		return nil, err
	}
	resp, err := a.Session.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a AtCoder) SqrapeQuizzes() ([]string, error) {
	resp, err := a.loginAndGet(a.QuizzesURL())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
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
	if len(quizzes) == 0 {
		return nil, fmt.Errorf("failed to scrape quiz from %v", a.QuizzesURL())
	}
	if invalidQuiz {
		return nil, fmt.Errorf("cannot get valid quizzes from %v", a.QuizzesURL())
	}
	return quizzes, nil
}

func (a AtCoder) SqrapeSample(quizID string) ([]*Sample, error) {
	resp, err := a.loginAndGet(a.QuizURL(quizID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	// AtCoder outputs both JP and EN content in one page.
	var ss []*Sample
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
	if len(ss) == 0 {
		return nil, fmt.Errorf("failed to scrape samples from %v", a.QuizURL(quizID))
	}
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
