package infra

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/garyburd/go-oauth/oauth"
)

func createWorkspace() error {
	name := configFilename()
	if _, err := os.Stat(name); err == nil {
		return nil
	}

	dir := filepath.Join(os.Getenv("HOME"), ".mint")
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	f, err := os.Create(name)
	if err != nil {
		return err
	}

	return f.Close()
}

func configFilename() string {
	return filepath.Join(os.Getenv("HOME"), ".mint/config.json")
}

func readJSON(src io.Reader, dest interface{}) error {
	return json.NewDecoder(src).Decode(dest)
}

type oauthRequest struct {
	cred        *oauth.Credentials
	method, url string
	params      url.Values
}
