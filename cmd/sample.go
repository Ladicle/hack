package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/Ladicle/hack/pkg"
	"github.com/spf13/cobra"
)

const FilePerm = 0644

type sampleCmd struct {
	TargetURL string
}

func NewSampleCmd() *cobra.Command {
	sample := sampleCmd{}
	cmd := &cobra.Command{
		Use:   "sample",
		Short: "Create sample files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sample.run(cmd, args)
		},
		SilenceUsage: true,
	}

	cmd.Flags().StringVar(&sample.TargetURL, "url", "",
		"[Required] Scraping target URL **NOTE: Supports AtCoder only**")
	cmd.MarkFlagRequired("url")

	return cmd
}

func (c *sampleCmd) run(cmd *cobra.Command, args []string) error {
	ss, err := pkg.SqrapeSample(c.TargetURL)
	if err != nil {
		return err
	}
	for _, sample := range ss {
		if err := ioutil.WriteFile(
			fmt.Sprintf("%v.in", sample.ID),
			[]byte(sample.Input),
			FilePerm,
		); err != nil {
			return err
		}
		if err := ioutil.WriteFile(
			fmt.Sprintf("%v.out", sample.ID),
			[]byte(sample.Output),
			FilePerm,
		); err != nil {
			return err
		}
		fmt.Printf("Create Sample #%v\n", sample.ID)
	}
	return nil
}
