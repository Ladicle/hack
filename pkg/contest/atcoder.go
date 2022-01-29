package contest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/Ladicle/hack/pkg/sample"
	"github.com/Ladicle/hack/pkg/session"
)

const (
	atCoderDir  = "atcoder"
	atCoderHost = "atcoder.jp"
)

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

// ScrapeContest scrapes contest data from /contests/<ContestID>/tasks.
func (a AtCoder) ScrapeContest() ([]string, error) {
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

// ScrapeSamples scrapes task samples from /contests/<ContestID>/tasks/<TaskID>.
func (a AtCoder) ScrapeSamples(taskID string) ([]*sample.Set, error) {
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

// SubmitCode submits the specified program.
func (a *AtCoder) SubmitCode(taskID, program string) error {
	code, err := ioutil.ReadFile(program)
	if err != nil {
		return err
	}
	ext := filepath.Ext(program)
	langId, err := ext2LangId(ext)
	if err != nil {
		return err
	}

	taskURL := GetTaskURL(a.ContestID, taskID)
	csrfToken, err := a.getCsrfToken(taskURL)
	if err != nil {
		return err
	}

	vals := url.Values{}
	vals.Add("data.TaskScreenName", taskID)
	vals.Add("data.LanguageId", langId)
	vals.Add("sourceCode", string(code))
	vals.Add("csrf_token", csrfToken)

	submitURL := fmt.Sprintf("https://%v/contests/%v/submit", atCoderHost, a.ContestID)
	resp, err := a.client.PostForm(submitURL, vals)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%v - %v", resp.StatusCode, resp.Status)
	}
	return nil
}

// ext2LangId returns the language ID from matched with the extension.
func ext2LangId(ext string) (string, error) {
	switch ext {
	case ".py":
		return "4047", nil
	case ".go":
		return "4026", nil
	case ".rs":
		return "4050", nil
	case ".cpp":
		return "4004", nil
	default:
		return "", fmt.Errorf("%q is not supported", ext)
	}
}

// GetTaskURL returns URL of the specified task page.
func GetTaskURL(contestID, taskID string) string {
	return fmt.Sprintf("https://%v/contests/%v/tasks/%v", atCoderHost, contestID, taskID)
}

// GetTaskURL returns my submission URL of the specified contest ID.
func GetSubmitMeURL(contestID string) string {
	return fmt.Sprintf("https://%v/contests/%v/submissions/me", atCoderHost, contestID)
}

// GetContestID returns the parent directory name as the contest ID.
func GetContestID(dir string) string {
	curBase := filepath.Base(dir)
	return filepath.Base(dir[:len(dir)-len(curBase)])
}

// GetTaskID returns the specified directory name as the task ID.
func GetTaskID(dir string) string {
	taskID := filepath.Base(dir)
	if strings.Contains(taskID, "_") {
		return taskID
	}
	contestID := GetContestID(dir)
	return fmt.Sprintf("%s_%s", contestID, filepath.Base(dir))
}

// GetAtCoderDir return the AtCoder directory name.
func GetAtCoderDir(baseDir string) string {
	return filepath.Join(baseDir, atCoderDir)
}

// GetContestDir return the contest directory name.
func GetContestDir(baseDir, contestID string) string {
	return filepath.Join(baseDir, atCoderDir, contestID)
}

// GetTaskDir return the task directory name.
func GetTaskDir(baseDir, contestID, taskID string) string {
	parts := strings.SplitN(taskID, "_", 2)
	if parts[0] == contestID {
		return filepath.Join(baseDir, parts[1])
	} else {
		return filepath.Join(baseDir, taskID)
	}
}
