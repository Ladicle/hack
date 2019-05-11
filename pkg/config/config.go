package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"github.com/spf13/viper"
)

const (
	DefaultConfig  = ".hack/config.yaml"
	DefaultBaseDir = "contest"

	baseDirKey        = "BaseDir"
	currentKey        = "Current"
	atCoderAccountKey = "AtCoderAccountKey"
	atCoderUserKey    = "AtCoderUser"
	atCoderPassKey    = "AtCoderPass"
	defaultLangKey    = "lang"
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
	glog.V(4).Infof("Use config file is %v", viper.ConfigFileUsed())
	baseDir := filepath.Dir(viper.ConfigFileUsed())
	glog.V(4).Infof("Create %v directory for the configuration", baseDir)
	if err := os.MkdirAll(baseDir, 0775); err != nil && !os.IsExist(err) {
		return err
	}
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
	return viper.GetString(fmt.Sprintf("%v.username", atCoderAccountKey))
}

func SetAtCoderAccount(account *Account) {
	viper.Set(atCoderAccountKey, account)
}

func AtCoderPass() string {
	return viper.GetString(fmt.Sprintf("%v.password", atCoderAccountKey))
}

func SetDefaultLang(lang string) {
	viper.Set(defaultLangKey, lang)
}

func DefaultLang() string {
	return viper.GetString(defaultLangKey)
}
