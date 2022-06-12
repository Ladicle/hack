package readme

import (
	"bytes"
	"text/template"
	"time"
)

const (
	orgTimeLayout = "[2006-01-02 Mon 15:04]"
	readmeTpl     = `{{define "base"}}#+title: {{ .ContestID }} - {{ .TaskID }}
#+setupfile: ~/doc/setup.prev.org
#+date: {{ .Date }}

* 問題

* 考察

* 解説

* 実装

#+include: "./main.py" src python
{{end}}
`
)

func GenerateReadme(contestID, taskID string, date time.Time) ([]byte, error) {
	tpl, err := template.New("base").Parse(readmeTpl)
	if err != nil {
		return nil, err
	}

	var data bytes.Buffer
	if err = tpl.Execute(&data, map[string]string{
		"ContestID": contestID,
		"TaskID":    taskID,
		"Date":      date.Format(orgTimeLayout),
	}); err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}
