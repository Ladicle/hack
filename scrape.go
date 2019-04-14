package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

const (
	targetURL = "https://atcoder.jp/contests/abs/tasks/abc085_b"
)

type Sample struct {
	ID     int
	Input  string
	Output string
}

func main() {
	ss, err := sqrapeSample(targetURL)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ss)
}

func sqrapeSample(target string) ([]*Sample, error) {
	var ss []*Sample

	doc, err := goquery.NewDocument(target)
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
				ID:    i,
				Input: s.Text(),
			})
		} else {
			ss[i/2-1].Output = s.Text()
		}
	})
	return ss, nil
}
