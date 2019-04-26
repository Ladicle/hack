package contest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	filePerm = 0644

	HostAtCoder = "atcoder"
)

type Sample struct {
	ID     int
	Input  string
	Output string
}

func mkQuizDir(contestDir string, quizzes []string) error {
	for _, n := range quizzes {
		if err := os.MkdirAll(filepath.Join(contestDir, n), 0755); err != nil {
			return err
		}
	}
	return nil
}

func mkSamples(quizDir string, samples []*Sample) error {
	for _, sample := range samples {
		if err := ioutil.WriteFile(
			filepath.Join(quizDir, fmt.Sprintf("%v.in", sample.ID)),
			[]byte(sample.Input),
			filePerm,
		); err != nil {
			return err
		}
		if err := ioutil.WriteFile(
			filepath.Join(quizDir, fmt.Sprintf("%v.out", sample.ID)),
			[]byte(sample.Output),
			filePerm,
		); err != nil {
			return err
		}
	}
	return nil
}
