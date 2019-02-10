package util

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	// LangGo represents the program is written in Go.
	LangGo = ".go"
	// LangCpp represents the program is written in C++.
	LangCpp = ".cpp"
)

// GetProgFName retuns a program file name in the specified directory.
func GetProgFName(dir string) (string, error) {
	var fname string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fname, err
	}

	for _, f := range files {
		if !f.IsDir() && strings.HasPrefix(f.Name(), "main") {
			fname = f.Name()
			break
		}
	}

	if fname == "" {
		return fname, fmt.Errorf("not found a program: the program must have 'main' prefix")
	}
	return fname, nil
}

// GetProgLang returns the specified program language.
func GetProgLang(fname string) (string, error) {
	if strings.HasSuffix(fname, LangGo) {
		return LangGo, nil
	} else if strings.HasSuffix(fname, LangCpp) {
		return LangCpp, nil
	}
	return "", fmt.Errorf("%v has unsupported extension", fname)
}
