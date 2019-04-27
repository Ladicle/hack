package config

import (
	"os/user"
	"path/filepath"
	"strings"

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

func CurrentContest() string {
	return viper.GetString(currentKey)
}

func CurrentContestPath() string {
	return filepath.Join(BaseDir(), viper.GetString(currentKey))
}

func CurrentHost() string {
	c := viper.GetString(currentKey)
	return strings.Split(c, "/")[0]
}

func CurrentContestID() string {
	c := viper.GetString(currentKey)
	return strings.TrimPrefix(strings.Split(c, "/")[1], "/")
}

func SetCurrentQuizPath(quiz string) string {
	return filepath.Join(CurrentContestPath(), quiz)
}

func AtCoderUser() string {
	return viper.GetString(atCoderUserKey)
}

func SetAtCoderUser(user string) {
	viper.Set(atCoderUserKey, user)
}

func AtCoderPass() string {
	return viper.GetString(atCoderPassKey)
}

func SetAtCoderPass(pass string) {
	viper.Set(atCoderPassKey, pass)
}
