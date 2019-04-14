package cmd

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestSampleCmd(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "hack_sample")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)

	cmd := sampleCmd{
		TargetURL: "https://atcoder.jp/contests/abs/tasks/abc085_b",
	}
	if err := cmd.run(&cobra.Command{}, []string{}); err != nil {
		t.Error(err)
	}

	finfo, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	expect := []string{"1.in", "1.out", "2.in", "2.out", "3.in", "3.out"}
	got := []string{}
	for _, f := range finfo {
		got = append(got, f.Name())
	}
	if !reflect.DeepEqual(expect, got) {
		t.Errorf("Generated unexpected files.\ngot: %v\nexpect: %v\n",
			expect, got)
	}
}
