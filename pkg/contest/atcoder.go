package contest

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/httputil"
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
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

func (a AtCoder) QuizzesURL() string {
	u, _ := url.Parse(atCoderBaseURL)
	u.Path = path.Join(u.Path, a.ContestID, atCoderQuizPath)
	return u.String()
}

func (a AtCoder) QuizURL(quizID string) string {
	u, _ := url.Parse(a.QuizzesURL())
	u.Path = path.Join(u.Path, quizID)
	return u.String()
}

func (a AtCoder) getCsrfToken(url string) (string, error) {
	glog.V(4).Infof("Getting csrfToken from %v...", url)
	resp, err := a.Session.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

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

func (a AtCoder) Login() error {
	if a.Session.Cookies != nil {
		glog.V(4).Info("Already logind to the AtCoder")
		return nil
	}

	glog.V(4).Info("Getting HTTP session...")
	csrfToken, err := a.getCsrfToken(atCoderLoginURL)
	if err != nil {
		return err
	}

	values := url.Values{}
	values.Add("username", config.AtCoderUser())
	values.Add("password", config.AtCoderPass())
	values.Add("csrf_token", csrfToken)

	glog.V(4).Infof("Login to the AtCoder as %q...", config.AtCoderUser())
	resp, err := a.Session.PostForm(atCoderLoginURL, &values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	glog.V(4).Info("Login succeeded")
	return nil
}

func (a AtCoder) loginAndGet(url string) (*http.Response, error) {
	if err := a.Login(); err != nil {
		return nil, err
	}
	glog.V(4).Infof("Getting %v...", url)
	resp, err := a.Session.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a AtCoder) ScrapeQuizzes() ([]string, error) {
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

func (a AtCoder) ScrapeSample(quizID string) ([]*Sample, error) {
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
