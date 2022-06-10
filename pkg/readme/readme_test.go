package readme

import "testing"

func TestGenerateReadme(t *testing.T) {
	const wantData = `#+title: sample-contest1 - sample-task1

* 問題

* 考察

* 解説

* 実装

#+include: "./main.py" src python
`

	gotData, err := GenerateReadme("sample-contest1", "sample-task1")
	if err != nil {
		t.Fatalf("fail to GenerateReadme(): err=%v", err)
	}
	if string(gotData) != wantData {
		t.Fatalf("got unexpected GenerateReadme() = %s, want=%s", gotData, wantData)
	}
}
