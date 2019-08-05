package infra

import (
	"os"
	"path/filepath"
)

func createWorkspace() error {
	dir := filepath.Join(os.Getenv("HOME"), ".mint")
	return os.MkdirAll(dir, 0700)
}

func configFilename() string {
	return filepath.Join(os.Getenv("HOME"), ".mint/config.json")
}
