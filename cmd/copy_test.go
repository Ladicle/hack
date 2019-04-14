package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Ladicle/hack/pkg/ttool"
)

func TestProgName(t *testing.T) {
	tmpDir := ttool.CreateAndGoWs(t)
	defer os.RemoveAll(tmpDir)

	_, err := getProgName()
	if err == nil {
		t.Error("Expected not found error")
	}

	f, err := ioutil.TempFile(tmpDir, "main.go")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(f.Name())
	f.Close()

	_, err = getProgName()
	if err != nil {
		t.Error(err)
	}
}
