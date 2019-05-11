package util

import (
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
	return "", nil
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

func InDir(name string) (bool, error) {
	fs, err := ioutil.ReadDir(".")
	if err != nil {
		return false, err
	}
	for _, f := range fs {
		if f.Name() == name {
			return true, nil
		}
	}
	return false, nil
}

// CleanRead read file and trim end of blank line.
func CleanRead(name string) (string, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(b), "\n"), nil
}
