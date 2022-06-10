package readme

import (
	"bytes"
	"text/template"
)

const readmeTpl = `{{define "base"}}#+title: {{ .ContestID }} - {{ .TaskID }}

* 問題

* 考察

* 解説

* 実装

#+include: "./main.py" src python
{{end}}
`

func GenerateReadme(contestID, taskID string) ([]byte, error) {
	tpl, err := template.New("base").Parse(readmeTpl)
	if err != nil {
		return nil, err
	}

	var data bytes.Buffer
	if err = tpl.Execute(&data, map[string]string{
		"ContestID": contestID,
		"TaskID":    taskID,
	}); err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}
