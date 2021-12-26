package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type File struct {
	AtCoder Account `yaml:"atcoder"`
}

type Account struct {
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

func MustUnmarshal(path string, cfg *File) {
	if err := Unmarshal(path, cfg); err != nil {
		log.Fatalf("Fail to load configuration file from %v: %v", path, err)
	}
}

func Unmarshal(path string, cfg *File) error {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = filepath.Join(home, path[1:])
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, cfg)
}
