package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/contest"
	"github.com/Ladicle/hack/pkg/pretty"
	"github.com/golang/glog"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
)

const FilePerm = 0644

type sampleCmd struct {
	progress *pretty.Progress
}

func NewSampleCmd() *cobra.Command {
	sample := sampleCmd{progress: pretty.NewProgress(os.Stdout)}
	cmd := &cobra.Command{
		Use:     "sample",
		Aliases: []string{"p"},
		Short:   "Create sample files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sample.Run(cmd, args)
		},
		SilenceUsage: true,
	}
	return cmd
}

func (c *sampleCmd) Run(cmd *cobra.Command, args []string) error {
	h := config.CurrentHost()
	if h == "" {
		return errors.New("contest is not set yet")
	}
	if h == contest.HostAtCoder {
		fmt.Printf("Creating %v samples...\n", h)
		return c.createAtCoderCurrentSamples()
	}
	// else
	// - get new sample id
	return nil
}

func (c *sampleCmd) createAtCoderCurrentSamples() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	quizID, err := currentQuizID(wd)
	if err != nil {
		return err
	}

	c.progress.Start(emoji.Sprintf("Scraping %v samples :inbox_tray:", quizID))
	at := contest.NewAtCoder(config.CurrentContestID())
	ss, err := at.ScrapeSample(quizID)
	if err != nil {
		c.progress.End(false)
		return err
	}

	if err := c.mkSamples(wd, ss); err != nil {
		return err
	}
	glog.V(4).Infof("Success to create AtCoder %v samples", quizID)
	return nil
}

func (c *sampleCmd) mkSamples(quizDir string, samples []*contest.Sample) error {
	for _, sample := range samples {
		c.progress.Start(emoji.Sprintf("Saving sample#%v :memo:", sample.ID))
		out := fmt.Sprintf("%v.in", sample.ID)
		if err := ioutil.WriteFile(
			filepath.Join(quizDir, out),
			[]byte(sample.Input),
			FilePerm,
		); err != nil {
			c.progress.End(false)
			return err
		}
		glog.V(4).Infof("Saved sample#%v to the %v", sample.ID, out)
		glog.V(8).Info(sample.Input)

		out = fmt.Sprintf("%v.out", sample.ID)
		if err := ioutil.WriteFile(
			filepath.Join(quizDir, out),
			[]byte(sample.Output),
			FilePerm,
		); err != nil {
			c.progress.End(false)
			return err
		}
		glog.V(4).Infof("Saved sample#%v to the %v", sample.ID, out)
		glog.V(8).Info(sample.Output)
	}
	c.progress.End(true)
	return nil
}

func currentQuizID(wd string) (string, error) {
	s := strings.TrimPrefix(wd, config.CurrentContestPath())
	if s == wd || s == "" {
		return "", fmt.Errorf("%q is not quiz directory", wd)
	}
	return s, nil
}
