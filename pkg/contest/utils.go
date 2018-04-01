package contest

import (
	"os"
	"path/filepath"
)

func mkdirs(baseDir string, dirNames []string) error {
	for _, n := range dirNames {
		if err := os.MkdirAll(filepath.Join(baseDir, n), 0755); err != nil {
			return err
		}
	}
	return nil
}
