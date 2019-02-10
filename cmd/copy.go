package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/Ladicle/hack/pkg/util"
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
	fname, err := util.GetProgFName(opt.WorkDir)
	if err != nil {
		return err
	}

	code, err := ioutil.ReadFile(filepath.Join(opt.WorkDir, fname))
	if err != nil {
		return err
	}

	if err := clipboard.WriteAll(string(code)); err != nil {
		return err
	}
	_, err = fmt.Fprintf(c.IO, "succeeded in copying %v to the clipboard\n", fname)
	return err
}
