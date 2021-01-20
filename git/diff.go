package git

import (
	"os/exec"
	"strings"
)

// GetChangedFiles returns a list of changed files
func GetChangedFiles() ([]string, error) {
	cmd := exec.Command(
		"git", "diff", "--name-only",
	)

	output, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}

	str := strings.TrimSpace(string(output))
	if str == "" {
		return []string{}, nil
	}
	return strings.Split(str, "\n"), nil
}

// GetUntrackedFiles returns a list of untracked files
func GetUntrackedFiles() ([]string, error) {
	cmd := exec.Command(
		"git", "ls-files", "--others", "--exclude-standard",
	)

	output, err := cmd.Output()
	if err != nil {
		return []string{}, err
	}

	str := strings.TrimSpace(string(output))
	if str == "" {
		return []string{}, nil
	}
	return strings.Split(str, "\n"), nil
}
