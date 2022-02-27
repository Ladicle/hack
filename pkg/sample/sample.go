package sample

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	ExtSampleIn  = ".in"
	ExtSampleOut = ".out"

	SampleDir = "samples"

	defaultEditor = "emacsclient"
)

type Set struct {
	In  string
	Out string
}

func (s Set) Write(dir string, id int, perm fs.FileMode) error {
	in := filepath.Join(dir, Name(id, ExtSampleIn))
	if err := ioutil.WriteFile(in, []byte(s.In), perm); err != nil {
		return err
	}
	out := filepath.Join(dir, Name(id, ExtSampleOut))
	if err := ioutil.WriteFile(out, []byte(s.Out), perm); err != nil {
		return err
	}
	return nil
}

func WriteInEditor(dir string, id int) error {
	for _, ext := range []string{ExtSampleIn, ExtSampleOut} {
		f := filepath.Join(dir, Name(id, ext))
		if err := exec.Command(defaultEditor, f).Run(); err != nil {
			return err
		}
	}
	return nil
}

// CntInputs counts the number of sample input files.
func CntInputs(dir string) (int, error) {
	entries, err := os.ReadDir(filepath.Join(dir, SampleDir))
	if err != nil {
		return -1, err
	}
	var cnt int
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ExtSampleIn {
			cnt++
		}
	}
	return cnt, nil
}

// Name returns the name of the sample input or output file with the specified id.
func Name(id int, ext string) string {
	return fmt.Sprintf("%s/%d%s", SampleDir, id, ext)
}
