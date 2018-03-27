package config

import (
	"io/ioutil"

	"os"

	"github.com/ghodss/yaml"
)

// C is singleton configuration data
var C Config

// LoadConfig loads configuration from path
func LoadConfig(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(f, &C)
}

// WriteConfig writes configuration to path
func WriteConfig(path string) error {
	y, err := yaml.Marshal(C)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, y, 0)
}
