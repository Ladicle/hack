package jump

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/contest"
	"github.com/Ladicle/hack/pkg/format"
	"github.com/Ladicle/hack/pkg/util"
	"github.com/golang/glog"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "jump [quiz]",
		Aliases: []string{"j"},
		Short:   "Get current quiz directory",
		Run: func(cmd *cobra.Command, args []string) {
			path := getPath(args)
			fmt.Println(path)
			if config.CurrentHost() == contest.HostAtCoder {
				base := filepath.Base(path)
				at, err := contest.NewAtCoder(config.CurrentContestID())
				if err != nil {
					glog.Fatal(err)
				}
				if err := open.Run(at.QuizURL(base)); err != nil {
					glog.Fatal(err)
				}
			}
		},
	}
}

func getPath(args []string) string {
	if len(args) >= 1 {
		return config.SetCurrentQuizPath(args[0])
	}

	hb := format.NewHackRobot(os.Stdout)
	cc := config.CurrentContestPath()
	fs, err := ioutil.ReadDir(cc)
	if err != nil {
		if os.IsNotExist(err) {
			hb.Fatal("Not exists the %q directory. Please set the contest.",
				config.CurrentContest())
		}
		glog.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		glog.Fatal(err)
	}

	quiz, err := contest.CurrentQuizID(wd)
	if err != nil {
		// In other directory
		quiz, err := firstQuizDir(fs)
		if err != nil {
			// No quiz directories are initialized
			return cc
		}
		return quiz
	}
	// In question directory
	dir, err := nextQuizDir(quiz, fs)
	if err != nil {
		hb.Fatal("%v", err)
	}
	return dir
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
	return "", fmt.Errorf("%q is unexpected quiz directory", currentQuiz)
}

func firstQuizDir(fsInDir []os.FileInfo) (string, error) {
	for _, f := range fsInDir {
		if !util.IsVisibleDir(f) {
			glog.V(8).Infof("Skip %q invisible directory", f.Name())
			continue
		}
		return config.SetCurrentQuizPath(f.Name()), nil
	}
	return "", errors.New("not in the current quiz directory")
}
