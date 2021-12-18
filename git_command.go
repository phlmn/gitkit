package gitkit

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

var gitCommandRegex = regexp.MustCompile(`^(git[-|\s]upload-pack|git[-|\s]upload-archive|git[-|\s]receive-pack) '(.*)'$`)

type GitCommand struct {
	Command  string
	Repo     string
	Original string
}

func ParseGitCommand(cmd string) (*GitCommand, error) {
	matches := gitCommandRegex.FindAllStringSubmatch(cmd, 1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("invalid git command")
	}

	// prevent path traversal
	safeRepo := path.Clean(path.Join("/", matches[0][2]))
	safeRepo = strings.TrimPrefix(safeRepo, "/")

	result := &GitCommand{
		Original: cmd,
		Command:  matches[0][1],
		Repo:     safeRepo,
	}

	return result, nil
}
