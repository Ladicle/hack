package jump

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/contest"
	"github.com/Ladicle/hack/pkg/util"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "jump [quiz]",
		Aliases: []string{"j"},
		Short:   "Get current quiz directory",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var quizPath string
			if len(args) >= 1 {
				quizPath = quizPathFromID(args[0])
			} else {
				quizPath, err = nextQuizPath()
			}
			fmt.Println(quizPath)
			return nil
		},
		SilenceUsage: true,
	}
}

func quizPathFromID(quizID string) string {
	return config.SetCurrentQuizPath(quizID)
}

func nextQuizPath() (string, error) {
	var nextQuiz string

	cc := config.CurrentContestPath()
	fs, err := ioutil.ReadDir(cc)
	if err != nil {
		return nextQuiz, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return nextQuiz, err
	}

	quiz, err := contest.CurrentQuizID(wd)
	if err != nil {
		// In other directory
		return firstQuizDir(fs)
	}
	// In question directory
	return nextQuizDir(quiz, fs)
}

func nextQuizDir(currentQuiz string, fsInDir []os.FileInfo) (string, error) {
	glog.V(4).Infof("Get next quiz directory: current=%v", currentQuiz)
	for i, f := range fsInDir {
		if !util.IsVisibleDir(f) {
			glog.V(8).Infof("Skip %q invisible directory", f.Name())
			continue
		}
		if f.Name() != currentQuiz {
			glog.V(8).Infof("Skip %q directory", f.Name())
			continue
		}
		if i+1 == len(fsInDir) {
			// There is no more quiz
			return config.SetCurrentQuizPath(fsInDir[i].Name()), nil
		}
		// return next quiz
		return config.SetCurrentQuizPath(fsInDir[i+1].Name()), nil
	}
	return "", errors.New("not found quiz directory")
}

func firstQuizDir(fsInDir []os.FileInfo) (string, error) {
	for _, f := range fsInDir {
		if !util.IsVisibleDir(f) {
			glog.V(8).Infof("Skip %q invisible directory", f.Name())
			continue
		}
		return config.SetCurrentQuizPath(f.Name()), nil
	}
	return "", errors.New("not found quiz directory")
}
