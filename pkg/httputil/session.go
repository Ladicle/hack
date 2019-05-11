package httputil

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/golang/glog"
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
	glog.V(4).Infof("Get %v", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return s.client.Do(req)
}

func (s *Session) PostForm(url string, v *url.Values) (*http.Response, error) {
	glog.V(4).Infof("Post %v to %v", v, url)
	req, err := http.NewRequest("POST", url, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return s.client.Do(req)
}
