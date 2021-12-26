package contest

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

const filePerm = 0644

type Sample struct {
	In  string
	Out string
}

func (s Sample) Write(dir string, id int) error {
	in := filepath.Join(dir, fmt.Sprintf("%d.in", id))
	if err := ioutil.WriteFile(in, []byte(s.In), filePerm); err != nil {
		return err
	}
	out := filepath.Join(dir, fmt.Sprintf("%d.out", id))
	if err := ioutil.WriteFile(out, []byte(s.Out), filePerm); err != nil {
		return err
	}
	return nil
}
