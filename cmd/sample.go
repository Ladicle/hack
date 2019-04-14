package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/Ladicle/hack/pkg"
	"github.com/spf13/cobra"
)

const FilePermission = 0644

var targetURL string

var sampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "Create sample files",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSample(cmd, args)
	},
	SilenceUsage: true,
}

func init() {
	sampleCmd.Flags().StringVar(&targetURL, "url", "",
		"[Required] Scraping target URL **NOTE: Supports AtCoder only**")
	sampleCmd.MarkFlagRequired("url")
}

func runSample(cmd *cobra.Command, args []string) error {
	ss, err := pkg.SqrapeSample(targetURL)
	if err != nil {
		return err
	}
	for _, sample := range ss {
		if err := ioutil.WriteFile(
			fmt.Sprintf("%v.in", sample.ID),
			[]byte(sample.Input),
			FilePermission,
		); err != nil {
			return err
		}
		if err := ioutil.WriteFile(
			fmt.Sprintf("%v.out", sample.ID),
			[]byte(sample.Output),
			FilePermission,
		); err != nil {
			return err
		}
		fmt.Printf("Create Sample #%v\n", sample.ID)
	}
	return nil
}
