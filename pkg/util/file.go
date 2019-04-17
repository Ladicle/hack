package util

import (
	"os"
	"strings"
)

func IsVisibleDir(f os.FileInfo) bool {
	return f.IsDir() && !strings.HasPrefix(f.Name(), ".")
}
