package pkg

import (
	"github.com/PuerkitoBio/goquery"
)

type Sample struct {
	ID     int
	Input  string
	Output string
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
