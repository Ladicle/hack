package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/Ladicle/hack/pkg/cmd/ini"
	"github.com/Ladicle/hack/pkg/cmd/open"
	"github.com/Ladicle/hack/pkg/cmd/test"
	"github.com/Ladicle/hack/pkg/config"
)

// set value in build time
var version string

const defaultPath = "~/.config/hack"

func Run() error {
	cmd := cobra.Command{
		Use:               "hack",
		Short:             "Hack assists your programming contest.",
		Version:           version,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	var (
		f   = &config.File{}
		out = os.Stdout

		path string
	)
	cmd.PersistentFlags().StringVar(&path, "config", defaultPath, "path to the configuration file")
	cmd.AddCommand(ini.NewCommand(f, out))
	cmd.AddCommand(test.NewCommand(f, out))
	cmd.AddCommand(open.NewCommand())

	cobra.OnInitialize(func() {
		config.MustUnmarshal(path, f)
	})
	return cmd.Execute()
}
