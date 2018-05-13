package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/atotto/clipboard"
)

func TestCopy(t *testing.T) {
	d, err := ioutil.TempDir("", "hack-copy")
	if err != nil {
		t.Fatalf("Could not create temporary directory: %v", err)
	}
	defer os.RemoveAll(d)

	want := "test data"
	ioutil.WriteFile(filepath.Join(d, "main.go"), []byte(want), 0644)

	c := copyCmd{IO: os.Stdout}
	if err := c.run([]string{}, Option{WorkDir: d}); err != nil {
		if clipboard.Unsupported {
			return
		}
		t.Fatalf("Failed to run copy command: %v", err)
	}

	if got, err := clipboard.ReadAll(); err != nil {
		t.Fatalf("Could not read clipboard data: %v", err)
	} else if got != want {
		t.Fatalf("cb data is different: want=%v, got=%v", want, got)
	}
}
