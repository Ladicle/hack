package cmd

import (
	"fmt"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/spf13/cobra"
)

var (
	gitRepo string
	version string

	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "hack",
	Short: "Hack assists your programming contest.",
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return config.Save()
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		fmt.Sprintf("config file (default %v)", config.DefaultCfg()))

	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(NewSampleCmd())
	rootCmd.AddCommand(NewCopyCmd())
	rootCmd.AddCommand(NewInitCmd())
}

func initConfig() {
	config.Load(cfgFile)
}

func Execute() error {
	return rootCmd.Execute()
}
