package cmd

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Ladicle/hack/pkg/config"
	"github.com/golang/glog"
	"golang.org/x/crypto/ssh/terminal"
)

func askBaseDir() (string, error) {
	var baseDir string
	fmt.Printf("# Where do you put contests code? (default: ~/%v)\n-> ",
		config.DefaultBaseDir)
	fmt.Scanln(&baseDir)
	fmt.Println()
	baseDir = strings.TrimSpace(baseDir)

	u, err := user.Current()
	if err != nil {
		return baseDir, err
	}

	if baseDir == "" {
		baseDir = filepath.Join(u.HomeDir, config.DefaultBaseDir)
	}
	if strings.HasPrefix(baseDir, "~") {
		baseDir = strings.Replace(baseDir, "~", u.HomeDir, 1)
	}
	glog.V(4).Info("Saved base directory")
	return baseDir, nil
}

func initAtCoder() (*config.Account, error) {
	var user, pass string
	fmt.Printf("## Tell me the username.\n-> ")
	fmt.Scanln(&user)
	fmt.Printf("## Tell me the password.\n-> ")
	bpass, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, err
	}
	pass = string(bpass)
	fmt.Printf("\n\n")

	ac := config.Account{
		Username: user,
		Password: pass,
	}
	glog.V(4).Info("Saved AtCoder username and password")
	return &ac, nil
}
