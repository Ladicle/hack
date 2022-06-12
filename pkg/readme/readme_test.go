package readme

import (
	"testing"
	"time"
)

func TestGenerateReadme(t *testing.T) {
	const wantData = `#+title: sample-contest1 - sample-task1
#+setupfile: ~/doc/setup.prev.org
#+date: [2022-06-12 Sun 21:38]

* 問題

* 考察

* 解説

* 実装

#+include: "./main.py" src python
`

	date, err := time.Parse(orgTimeLayout, "[2022-06-12 Sun 21:38]")
	if err != nil {
		t.Fatalf("fail to parse time: err=%v", err)
	}
	gotData, err := GenerateReadme("sample-contest1", "sample-task1", date)
	if err != nil {
		t.Fatalf("fail to GenerateReadme(): err=%v", err)
	}
	if string(gotData) != wantData {
		t.Fatalf("got unexpected GenerateReadme() = %s, want=%s", gotData, wantData)
	}
}
