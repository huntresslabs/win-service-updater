package updater

import (
	"os"
	"path/filepath"
)

func GetExeDir() string {
	exe, _ := os.Executable()
	return filepath.Dir(exe)
}
