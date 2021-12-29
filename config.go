package gitkit

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	KeyDir     string       // Directory for server ssh keys. Only used in SSH strategy.
	Dir        string       // Directory that contains repositories
	GitPath    string       // Path to git binary
	GitUser    string       // User for ssh connections
	AutoCreate bool         // Automatically create repostories
	AutoHooks  bool         // Automatically setup git hooks
	Hooks      *HookScripts // Scripts for hooks/* directory
	Auth       bool         // Require authentication
}

// HookScripts represents all repository server-size git hooks
type HookScripts struct {
	PreReceive  string
	Update      string
	PostReceive string
}

// Configure hook scripts in the repo base directory
func (c *HookScripts) setupInDir(path string) error {
	basePath := filepath.Join(path, "hooks")
	scripts := map[string]string{
		"pre-receive":  c.PreReceive,
		"update":       c.Update,
		"post-receive": c.PostReceive,
	}

	// Cleanup any existing hooks first
	hookFiles, err := ioutil.ReadDir(basePath)
	if err == nil {
		for _, file := range hookFiles {
			if err := os.Remove(filepath.Join(basePath, file.Name())); err != nil {
				return err
			}
		}
	}

	// Write new hook files
	for name, script := range scripts {
		fullPath := filepath.Join(basePath, name)

		// Dont create hook if there's no script content
		if script == "" {
			continue
		}

		if err := ioutil.WriteFile(fullPath, []byte(script), 0755); err != nil {
			logError("hook-update", err)
			return err
		}
	}

	return nil
}

func (c *Config) KeyPath(keyType string) string {
	return filepath.Join(c.KeyDir, "gitkit"+"."+keyType)
}

func (c *Config) Setup() error {
	if _, err := os.Stat(c.Dir); err != nil {
		if err = os.Mkdir(c.Dir, 0755); err != nil {
			return err
		}
	}

	if c.AutoHooks == true {
		return c.setupHooks()
	}

	return nil
}

func (c *Config) setupHooks() error {
	walk := func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && strings.HasSuffix(d.Name(), ".git") {
			if err := c.Hooks.setupInDir(s); err != nil {
				return err
			}
		}
		return nil
	}

	if err := filepath.WalkDir(c.Dir, walk); err != nil {
		return err
	}

	return nil
}
