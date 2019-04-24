package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	defaultConfig  = ".hack.yaml"
	defaultBaseDir = "contest"

	baseDirKey     = "BaseDir"
	currentKey     = "Current"
	atCoderUserKey = "AtCoderUser"
	atCoderPassKey = "AtCoderPass"
)

func Load(overwriteCfg string) {
	viper.SetConfigFile(DefaultCfg())
	if overwriteCfg != "" {
		viper.SetConfigFile(overwriteCfg)
	}

	if err := viper.ReadInConfig(); err != nil {
		if !os.IsNotExist(err) {
			fmt.Println(err)
			os.Exit(1)
		}
		// set defaults at first time
		setDefaults()
	}
}

func setDefaults() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	viper.SetDefault(baseDirKey, filepath.Join(u.HomeDir, defaultBaseDir))
	return nil
}

func Save() error {
	return viper.WriteConfig()
}

func DefaultCfg() string {
	u, err := user.Current()
	if err != nil {
		os.Exit(1)
	}
	return filepath.Join(u.HomeDir, defaultConfig)
}

func SetCurrent(c string) {
	viper.Set(currentKey, c)
}

func BaseDir() string {
	return viper.GetString(baseDirKey)
}

func CurrentContestPath() string {
	return filepath.Join(BaseDir(), viper.GetString(currentKey))
}

func CurrentQuizPath(quiz string) string {
	return filepath.Join(CurrentContestPath(), quiz)
}

func AtCoderUser() string {
	return viper.GetString(atCoderUserKey)
}

func AtCoderPass() string {
	return viper.GetString(atCoderPassKey)
}
