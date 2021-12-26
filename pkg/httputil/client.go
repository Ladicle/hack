package httputil

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type Session struct {
	client *http.Client
}

func NewSession(serverName string) (*Session, error) {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, err
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			ServerName: serverName,
		},
	}
	client := &http.Client{
		Jar:       jar,
		Transport: tr,
	}
	return &Session{client: client}, nil
}

func (s *Session) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: [GET] %s", err, url)
	}
	resp, err := do(s.client, req)
	if err != nil {
		return nil, fmt.Errorf("%w: [GET] %s", err, url)
	}
	return resp, nil
}

func (s *Session) PostForm(url string, v *url.Values) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, fmt.Errorf("%w: [POST] %s?%s", err, url, v.Encode())
	}
	resp, err := do(s.client, req)
	if err != nil {
		return nil, fmt.Errorf("%w: [POST] %s?%s", err, url, v.Encode())
	}
	return resp, nil
}

func do(c *http.Client, req *http.Request) (*http.Response, error) {
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	return resp, nil
}
