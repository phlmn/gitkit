package gitkit

import (
	"os"
	"os/exec"
	"path"
)

func InitRepo(name string, config *Config) error {
	fullPath := path.Join(config.Dir, name)

	if err := exec.Command(config.GitPath, "init", "--bare", fullPath).Run(); err != nil {
		return err
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
