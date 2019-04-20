package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func IsVisibleDir(f os.FileInfo) bool {
	return f.IsDir() && !strings.HasPrefix(f.Name(), ".")
}

func GetProgName() (string, error) {
	var fname string

	fs, err := ioutil.ReadDir(".")
	if err != nil {
		return fname, err
	}

	for _, f := range fs {
		if !f.IsDir() && strings.HasPrefix(f.Name(), "main.") {
			return f.Name(), nil
		}
	}

	return fname, fmt.Errorf("not found a main program")
}

func SampleIDs(dir string) ([]string, error) {
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var samples []string
	for _, f := range fs {
		if !strings.HasSuffix(f.Name(), ".out") {
			continue
		}
		id := strings.SplitN(f.Name(), ".", 2)[0]
		samples = append(samples, id)
	}
	return samples, nil
}
