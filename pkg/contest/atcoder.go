package contest

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/httputil"
	"github.com/Ladicle/hack/pkg/lang"
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
)

const (
	atCoderBaseURL  = "https://atcoder.jp/contests"
	atCoderLoginURL = "https://atcoder.jp/login"

	atCoderQuizPath      = "tasks"
	atCoderSubmitPath    = "submit"
	atCoderSubmittedPath = "submissions/me"
)

func NewAtCoder(id string) (*AtCoder, error) {
	s, err := httputil.NewSession("atcoder.jp")
	if err != nil {
		return nil, err
	}
	return &AtCoder{
		ContestID: id,
		session:   s,
	}, nil
}

type AtCoder struct {
	ContestID string

	session   *httputil.Session
	csrfToken string
}

func (a *AtCoder) QuizzesURL() string {
	u, _ := url.Parse(atCoderBaseURL)
	u.Path = path.Join(u.Path, a.ContestID, atCoderQuizPath)
	return u.String()
}

func (a *AtCoder) QuizURL(quizID string) string {
	u, _ := url.Parse(a.QuizzesURL())
	u.Path = path.Join(u.Path, quizID)
	return u.String()
}

func (a *AtCoder) SubmitURL() string {
	u, _ := url.Parse(atCoderBaseURL)
	u.Path = path.Join(u.Path, a.ContestID, atCoderSubmitPath)
	return u.String()
}

func (a *AtCoder) SubmittedURL() string {
	u, _ := url.Parse(atCoderBaseURL)
	u.Path = path.Join(u.Path, a.ContestID, atCoderSubmittedPath)
	return u.String()
}

func (a *AtCoder) getCsrfToken(url string) (string, error) {
	glog.V(4).Infof("Getting csrfToken from %v...", url)
	resp, err := a.session.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%v - %v", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", err
	}

	glog.V(4).Info("Scrapping csrfToken from response...")
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
	glog.V(4).Info("Success to get csrfToken")
	return csrfToken, nil
}

func (a *AtCoder) Login() error {
	glog.V(4).Info("Getting HTTP session...")
	csrfToken, err := a.getCsrfToken(atCoderLoginURL)
	if err != nil {
		return err
	}
	a.csrfToken = csrfToken

	values := url.Values{}
	values.Add("username", config.AtCoderUser())
	values.Add("password", config.AtCoderPass())
	values.Add("csrf_token", csrfToken)

	glog.V(4).Infof("Login to the AtCoder as %q...", config.AtCoderUser())
	resp, err := a.session.PostForm(atCoderLoginURL, &values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%v - %v", resp.StatusCode, resp.Status)
	}

	glog.V(4).Info("Login succeeded")
	return nil
}

func (a *AtCoder) loginAndGet(url string) (*http.Response, error) {
	if err := a.Login(); err != nil {
		return nil, err
	}
	glog.V(4).Infof("Getting %v...", url)
	resp, err := a.session.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("%v - %v", resp.StatusCode, resp.Status)
	}
	return resp, nil
}

func (a *AtCoder) loginAndPost(URL string, values *url.Values) (*http.Response, error) {
	if err := a.Login(); err != nil {
		return nil, err
	}

	csrfToken, err := a.getCsrfToken(URL)
	if err != nil {
		return nil, err
	}
	if csrfToken == "" {
		return nil, errors.New("csrfToken is null")
	}
	values.Add("csrf_token", csrfToken)

	glog.V(4).Infof("Posting %v...", URL)
	resp, err := a.session.PostForm(URL, values)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("%v - %v", resp.StatusCode, resp.Status)
	}
	return resp, nil
}

func (a *AtCoder) ScrapeQuizzes() ([]string, error) {
	resp, err := a.loginAndGet(a.QuizzesURL())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	glog.V(4).Infof("Scrapping quizzes from %v...", a.QuizzesURL())
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	var quizzes []string
	var invalidQuiz bool
	doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
		path, ok := s.Find("a").First().Attr("href")
		if !ok {
			glog.Warningf("Node #%v: has not URL attribute", i)
			invalidQuiz = true
		}
		// path has /contests/<id>/tasks/<quiz_id>
		section := strings.Split(path, "/")
		quizzes = append(quizzes, section[len(section)-1])
		glog.V(8).Infof("Node #%v: save quizID from %v", i, section)
	})
	if len(quizzes) == 0 {
		return nil, fmt.Errorf("failed to scrape quiz from %v", a.QuizzesURL())
	}
	if invalidQuiz {
		return nil, fmt.Errorf("cannot get valid quizzes from %v", a.QuizzesURL())
	}
	glog.V(4).Info("Success to scrape quizzes")
	return quizzes, nil
}

func (a *AtCoder) ScrapeSample(quizID string) ([]*Sample, error) {
	resp, err := a.loginAndGet(a.QuizURL(quizID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	glog.V(4).Infof("Scraping samples from %v...", a.QuizURL(quizID))
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	// AtCoder outputs both JP and EN content in one page.
	var ss []*Sample
	doc.Find(".lang-en pre").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			glog.V(8).Infof("Node #%v: skip it because first block is not sample", i)
			return
		}
		if i%2 == 1 {
			id := i/2 + 1
			glog.V(8).Infof("Node #%v: append a sample as %v.in", i, id)
			ss = append(ss, &Sample{ID: id, Input: s.Text()})
			return
		}
		glog.V(8).Infof("Node #%v: append a sample as %v.out", i, i/2)
		ss[i/2-1].Output = s.Text()
	})
	if len(ss) == 0 {
		return nil, fmt.Errorf("failed to scrape samples from %v", a.QuizURL(quizID))
	}
	glog.V(4).Info("Success to scrape samples")
	return ss, nil
}

func (a *AtCoder) SubmitCode(quizID, sorceFile string) error {
	code, err := ioutil.ReadFile(sorceFile)
	if err != nil {
		return err
	}
	ext := filepath.Ext(sorceFile)
	langId, err := ext2LangId(ext)
	if err != nil {
		return err
	}
	values := url.Values{}
	values.Add("data.TaskScreenName", quizID)
	values.Add("data.LanguageId", langId)
	values.Add("sourceCode", string(code))

	glog.V(4).Infof("Submit %q solution to the AtCoder...", quizID)
	resp, err := a.loginAndPost(a.SubmitURL(), &values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	glog.V(4).Infof("Success to submit code: %v", buf.String())

	return nil
}

func ext2LangId(ext string) (string, error) {
	switch ext {
	case lang.TypeCpp:
		return "3003", nil
	case lang.TypeGo:
		return "3013", nil
	}
	return "", fmt.Errorf("%q is not supported", ext)
}
