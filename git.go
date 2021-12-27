package gitkit

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
)

func InitRepo(name string, config *Config) error {
	fullPath := path.Join(config.Dir, name)

	cmd := exec.Command(config.GitPath, "init", "--bare", fullPath)
	stderr, _ := cmd.StdoutPipe()
	stderrStr, _ := io.ReadAll(stderr)

	if err := cmd.Run(); err != nil {
		return errors.New(fmt.Sprintf("%s: %s", err, stderrStr))
	}

	if config.AutoHooks && config.Hooks != nil {
		return config.Hooks.setupInDir(fullPath)
	}

	return nil
}

func CloneRepo(name string, config *Config, url string) error {
	fullPath := path.Join(config.Dir, name)

	cmd := exec.Command(config.GitPath, "clone", "--bare", url, fullPath)
	stderr, _ := cmd.StdoutPipe()
	stderrStr, _ := io.ReadAll(stderr)

	if err := cmd.Run(); err != nil {
		return errors.New(fmt.Sprintf("%s: %s", err, stderrStr))
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
