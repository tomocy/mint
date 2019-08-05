package infra

import (
	"os"
	"path/filepath"
)

func configFilename() string {
	return filepath.Join(os.Getenv("HOME"), ".mint/config.json")
}
