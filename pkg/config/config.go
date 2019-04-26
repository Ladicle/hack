package config

import (
	"os/user"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	DefaultConfig  = ".hack.yaml"
	DefaultBaseDir = "contest"

	baseDirKey     = "BaseDir"
	currentKey     = "Current"
	atCoderUserKey = "AtCoderUser"
	atCoderPassKey = "AtCoderPass"
)

func Load(overwriteCfg string) error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	viper.SetConfigFile(filepath.Join(u.HomeDir, DefaultConfig))
	if overwriteCfg != "" {
		viper.SetConfigFile(overwriteCfg)
	}
	return viper.ReadInConfig()
}

func setDefaults() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	viper.SetDefault(baseDirKey, filepath.Join(u.HomeDir, DefaultBaseDir))
	return nil
}

func Save() error {
	return viper.WriteConfig()
}

func SetCurrent(c string) {
	viper.Set(currentKey, c)
}

func SetBaseDir(dir string) {
	viper.Set(baseDirKey, dir)
}

func BaseDir() string {
	return viper.GetString(baseDirKey)
}

func CurrentContestPath() string {
	return filepath.Join(BaseDir(), viper.GetString(currentKey))
}

func SetCurrentQuizPath(quiz string) string {
	return filepath.Join(CurrentContestPath(), quiz)
}

func AtCoderUser() string {
	return viper.GetString(atCoderUserKey)
}

func AtCoderPass() string {
	return viper.GetString(atCoderPassKey)
}
