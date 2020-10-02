package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Branch defines a git branch
type Branch struct {
	Name          string
	LastCommit    time.Time
	LastCommitter string
}

// fields on Branch struct to display
var fieldsToDisplay = []string{"Name", "LastEdited", "LastCommitter"}

// fetch branch metadata (last commit, committer)
func (b *Branch) populateBranchMetadata() {
	gitExecutable, _ := exec.LookPath("git")
	cmd := exec.Command(
		gitExecutable,
		"--no-pager",
		"log",
		"-1",
		// https://git-scm.com/docs/pretty-formats
		`--pretty=format:%ct %cn`,
		b.Name,
	)

	output, err := cmd.Output()
	if err != nil {
		return
	}

	fields := strings.Fields(string(output))
	if len(fields) == 2 {
		b.LastCommit = parseUnixTimestamp(fields[0])
		b.LastCommitter = fields[1]
		fmt.Println(b.formatDate())
	}
}

func parseUnixTimestamp(str string) time.Time {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.Unix(i, 0)
}

func (b *Branch) formatDate() string {
	if b.LastCommit.IsZero() {
		return ""
	}
	year, month, day := b.LastCommit.Date()
	return fmt.Sprintf("%v-%02v-%02v", year, int(month), day)
}
