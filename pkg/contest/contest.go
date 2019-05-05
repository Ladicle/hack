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
	dirPerm  = 0775

	HostAtCoder = "atcoder"
	HostCodeJam = "codejam"
	HostFree    = "free"
)

type Sample struct {
	ID     int
	Input  string
	Output string
}

func MkCurrentContestDir() error {
	return os.MkdirAll(config.CurrentContestPath(), dirPerm)
}

func MkQuizDir(quizzes []string) error {
	for _, n := range quizzes {
		if err := os.Mkdir(n, dirPerm); err != nil {
			return err
		}
	}
	return nil
}

func MkSample(quizDir string, sample *Sample) error {
	out := fmt.Sprintf("%v.in", sample.ID)
	if err := ioutil.WriteFile(
		filepath.Join(quizDir, out),
		[]byte(sample.Input),
		filePerm,
	); err != nil {
		return err
	}
	glog.V(4).Infof("Saved sample#%v to the %v", sample.ID, out)
	glog.V(8).Info(sample.Input)

	out = fmt.Sprintf("%v.out", sample.ID)
	if err := ioutil.WriteFile(
		filepath.Join(quizDir, out),
		[]byte(sample.Output),
		filePerm,
	); err != nil {
		return err
	}
	glog.V(4).Infof("Saved sample#%v to the %v", sample.ID, out)
	glog.V(8).Info(sample.Output)
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
