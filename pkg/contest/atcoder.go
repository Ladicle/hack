package contest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/Ladicle/hack/pkg/httputil"
	"github.com/Ladicle/hack/pkg/lang"
)

const atCoderHost = "atcoder.jp"

type AtCoder struct {
	ContestID string

	client *httputil.Session
}

// NewAtCoder creates and return the AtCoder object.
func NewAtCoder(contestID string) (*AtCoder, error) {
	s, err := httputil.NewSession("atcoder.jp")
	if err != nil {
		return nil, err
	}
	return &AtCoder{
		ContestID: contestID,
		client:    s,
	}, nil
}

func (a *AtCoder) Login(user, pass string) error {
	addr := fmt.Sprintf("https://%v/login", atCoderHost)
	csrfToken, err := a.getCsrfToken(addr)
	if err != nil {
		return err
	}

	values := url.Values{}
	values.Add("username", user)
	values.Add("password", pass)
	values.Add("csrf_token", csrfToken)

	resp, err := a.client.PostForm(addr, &values)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

func (a *AtCoder) getCsrfToken(addr string) (string, error) {
	resp, err := a.client.Get(addr)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	var csrfToken string
	doc.Find("input[name=\"csrf_token\"]").Each(func(i int, s *goquery.Selection) {
		if value, exist := s.Attr("value"); exist {
			csrfToken = value
		}
	})
	if csrfToken == "" {
		return "", fmt.Errorf("failed to scrape csrfToken from %#+v", doc)
	}
	return csrfToken, nil
}

func (a *AtCoder) ScrapeTasks() ([]string, error) {
	addr := fmt.Sprintf("https://%v/contests/%v/tasks", atCoderHost, a.ContestID)
	resp, err := a.client.Get(addr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var (
		cnt   int
		tasks []string
	)
	doc.Find("table tbody tr td").Each(func(i int, s *goquery.Selection) {
		path, ok := s.Find("a").First().Attr("href")
		if ok && cnt%2 == 0 {
			// path has /contests/<id>/tasks/<quiz_id (e.g. abc228_a)>
			parts := strings.Split(path, "/")
			tasks = append(tasks, parts[len(parts)-1])
		}
		cnt++
	})
	if len(tasks) == 0 {
		return nil, fmt.Errorf("failed to scrape quiz from %#+v", doc)
	}
	return tasks, nil
}

func (a *AtCoder) ScrapeTask(taskID string) ([]*Sample, error) {
	addr := fmt.Sprintf("https://%v/contests/%v/tasks/%v", atCoderHost, a.ContestID, taskID)
	resp, err := a.client.Get(addr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	// AtCoder outputs both JP and EN content in one page.
	var ss []*Sample
	doc.Find(".lang-en pre").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		if i%2 == 1 {
			ss = append(ss, &Sample{In: s.Text()})
			return
		}
		ss[i/2-1].Out = s.Text()
	})
	if len(ss) == 0 {
		return nil, fmt.Errorf("failed to scrape samples from %#+v", doc)
	}
	return ss, nil
}

func (a *AtCoder) Submit(taskID, prog string) error {
	code, err := ioutil.ReadFile(prog)
	if err != nil {
		return err
	}
	ext := filepath.Ext(prog)
	langID, err := ext2LangID(ext)
	if err != nil {
		return err
	}

	values := url.Values{}
	values.Add("data.TaskScreenName", taskID)
	values.Add("data.LanguageId", langID)
	values.Add("sourceCode", string(code))

	addr := fmt.Sprintf("https://%v/contests/%v/submit", atCoderHost, a.ContestID)
	resp, err := a.client.PostForm(addr, &values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return err
	}
	return nil
}

func ext2LangID(ext string) (string, error) {
	switch ext {
	case lang.LangCpp:
		return "4003", nil
	case lang.LangGo:
		return "4026", nil
	case lang.LangPython:
		return "4006", nil
	}
	return "", fmt.Errorf("%q is unsupported language", ext)
}
