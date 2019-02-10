package cmd

import (
	"fmt"
	"io"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/ghodss/yaml"
)

// NewInfoCmd infos contest information.
func NewInfoCmd(io io.Writer) Command {
	s := infoCmd{IO: io}
	return Command{
		Name:        "info",
		Short:       "info",
		Description: "Info shows information",
		Run:         s.run,
	}
}

type infoCmd struct {
	IO io.Writer
}

func (c *infoCmd) run(args []string, opt Option) error {
	y, err := yaml.Marshal(config.C)
	if err != nil {
		return fmt.Errorf("failed to convert information to YAML: %v", err)
	}
	_, err = fmt.Fprintf(c.IO, string(y))
	return err
}
