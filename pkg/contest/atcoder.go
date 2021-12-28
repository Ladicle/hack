package contest

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/Ladicle/hack/pkg/sample"
	"github.com/Ladicle/hack/pkg/session"
)

const atCoderHost = "atcoder.jp"

type AtCoder struct {
	ContestID string

	client *session.Client
}

// NewAtCoder creates and return the AtCoder object.
func NewAtCoder(contestID string) (*AtCoder, error) {
	s, err := session.NewClient()
	if err != nil {
		return nil, err
	}
	return &AtCoder{
		ContestID: contestID,
		client:    s,
	}, nil
}

func (a AtCoder) Login(user, pass string) error {
	addr := fmt.Sprintf("https://%s/login", atCoderHost)
	csrfToken, err := a.getCsrfToken(addr)
	if err != nil {
		return err
	}

	vals := make(url.Values, 3)
	vals.Add("username", user)
	vals.Add("password", pass)
	vals.Add("csrf_token", csrfToken)

	resp, err := a.client.PostForm(addr, vals)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

func (a AtCoder) getCsrfToken(addr string) (string, error) {
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

func (a AtCoder) ScrapeTasks() ([]string, error) {
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

	var tasks []string
	doc.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
		path, ok := s.Find("td").First().Find("a").First().Attr("href")
		if ok {
			// path has /contests/<id>/tasks/<quiz_id (e.g. abc228_a)>
			parts := strings.Split(path, "/")
			tasks = append(tasks, parts[len(parts)-1])
		}
	})
	if len(tasks) == 0 {
		return nil, fmt.Errorf("failed to scrape quiz from %#+v", doc)
	}
	return tasks, nil
}

func (a AtCoder) ScrapeTask(taskID string) ([]*sample.Set, error) {
	addr := GetTaskURL(a.ContestID, taskID)
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
	var ss []*sample.Set
	doc.Find(".lang-en pre").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		if i%2 == 1 {
			ss = append(ss, &sample.Set{In: s.Text()})
			return
		}
		ss[i/2-1].Out = s.Text()
	})
	if len(ss) == 0 {
		return nil, fmt.Errorf("failed to scrape samples from %#+v", doc)
	}
	return ss, nil
}

// GetTaskURL returns URL of the specified task page.
func GetTaskURL(contestID, taskID string) string {
	return fmt.Sprintf("https://%v/contests/%v/tasks/%v", atCoderHost, contestID, taskID)
}

// GetContestID returns the parent directory name as the contest ID.
func GetContestID(dir string) string {
	curBase := filepath.Base(dir)
	return filepath.Base(dir[:len(dir)-len(curBase)])
}

// GetTaskID returns the specified directory name as the task ID.
func GetTaskID(dir string) string {
	return filepath.Base(dir)
}
