package contest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/golang/glog"
)

const (
	filePerm = 0644
	dirPerm  = 0755

	HostAtCoder = "atcoder"
)

type Sample struct {
	ID     int
	Input  string
	Output string
}

func MkCurrentContestDir() error {
	return os.MkdirAll(config.CurrentContestPath(), dirPerm)
}

func mkQuizDir(contestDir string, quizzes []string) error {
	for _, n := range quizzes {
		if err := os.MkdirAll(filepath.Join(contestDir, n), dirPerm); err != nil {
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

func CurrentQuizID(wd string) (string, error) {
	s := strings.TrimPrefix(wd, config.CurrentContestPath())
	if s == wd || s == "" {
		return s, fmt.Errorf("%q is not quiz directory", wd)
	}
	glog.V(8).Infof("CurrentQuiz: %q", s)
	return strings.TrimPrefix(s, "/"), nil
}
