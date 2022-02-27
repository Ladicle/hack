package add

import (
	"fmt"
	"io"
	"os"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/Ladicle/hack/pkg/sample"
	"github.com/spf13/cobra"
)

const resourceSample = "sample"

type Options struct {
	resource string
}

func NewCommand(f *config.File, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:          "add [<RESOURCE>]",
		Aliases:      []string{"a"},
		Short:        "Add resource to the contest directory",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var opts Options
			if err := opts.Validate(args); err != nil {
				return err
			}
			return opts.Run(f, out)
		},
	}
}

func (o *Options) Validate(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("invalid number of arguments: got=%v, want=%v", len(args), 1)
	}
	switch resource := args[0]; resource {
	case resourceSample:
		o.resource = resource
	default:
		return fmt.Errorf("unexpected resource: got=%v", resource)
	}
	return nil
}

func (o Options) Run(f *config.File, out io.Writer) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	cur, err := sample.CntInputs(wd)
	if err != nil {
		return err
	}
	return sample.WriteInEditor(wd, cur+1)
}
