package httputil

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Session struct {
	Cookies []*http.Cookie
}

func (s *Session) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for _, c := range s.Cookies {
		req.AddCookie(c)
	}
	return req, nil
}

func (s *Session) Do(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	s.Cookies = resp.Cookies()
	return resp, nil
}

func (s *Session) Get(url string) (*http.Response, error) {
	req, err := s.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return s.Do(req)
}

func (s *Session) PostForm(url string, v *url.Values) (*http.Response, error) {
	req, err := s.NewRequest("POST", url, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return s.Do(req)
}
