package gitkit

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
)

func InitRepo(name string, config *Config) error {
	fullPath := path.Join(config.Dir, name)

	var stderr string
	cmd := exec.Command(config.GitPath, "init", "--bare", fullPath)
	cmd.Stderr = bytes.NewBufferString(stderr)

	if err := cmd.Run(); err != nil {
		return errors.New(fmt.Sprintf("%s: %s", err, stderr))
	}

	if config.AutoHooks && config.Hooks != nil {
		return config.Hooks.setupInDir(fullPath)
	}

	return nil
}

func CloneRepo(name string, config *Config, url string) error {
	fullPath := path.Join(config.Dir, name)

	var stderr string
	cmd := exec.Command(config.GitPath, "clone", "--bare", url, fullPath)
	cmd.Stderr = bytes.NewBufferString(stderr)

	if err := cmd.Run(); err != nil {
		return errors.New(fmt.Sprintf("%s: %s", err, stderr))
	}

	if config.AutoHooks && config.Hooks != nil {
		return config.Hooks.setupInDir(fullPath)
	}

	return nil
}

func RepoExists(p string) bool {
	_, err := os.Stat(path.Join(p, "objects"))
	return err == nil
}
