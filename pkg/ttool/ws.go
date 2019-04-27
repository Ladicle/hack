package ttool

import (
	"io/ioutil"
	"os"
	"testing"
)

func CreateAndGoWs(t *testing.T) string {
	dir, err := ioutil.TempDir("", "hack")
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(dir); err != nil {
		os.RemoveAll(dir)
		t.Fatal(err)
	}
	return dir
}
