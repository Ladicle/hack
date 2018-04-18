package cmd

import (
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/atotto/clipboard"
)

// NewCopyCmd copies code to clipboard.
func NewCopyCmd(io io.Writer) Command {
	s := copyCmd{IO: io}
	return Command{
		Name:        "copy",
		Short:       "copy",
		Description: "Copy copies your code to clipboard",
		Run:         s.run,
	}
}

type copyCmd struct {
	IO io.Writer
}

func (c *copyCmd) run(args []string, opt Option) error {
	code, err := ioutil.ReadFile(filepath.Join(opt.WorkDir, "main.go"))
	if err != nil {
		return err
	}
	return clipboard.WriteAll(string(code))
}
