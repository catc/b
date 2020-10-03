package git

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

const dateFieldWidth = 11
const committerFieldPadding = 5
const nameFieldPadding = 1

// Branches contains all branches and git related configs
type Branches struct {
	Branches           []Branch
	CurrentBranch      string
	CurrentBranchIndex int
	MaxColumnWidth     map[string]int
}

func newGitBranches() *Branches {
	return &Branches{
		Branches:           make([]Branch, 0),
		CurrentBranch:      "_",
		CurrentBranchIndex: -1,
		MaxColumnWidth:     make(map[string]int),
	}
}

// GetBranches returns struct for git branches in the repo
func GetBranches() (*Branches, error) {
	gitExecutable, err := exec.LookPath("git")
	if err != nil {
		panic("git not found")
	}

	/*
		format:
		1601742631 author master *
	*/
	cmd := exec.Command(
		gitExecutable,
		"for-each-ref",
		"--sort=-authordate:iso8601",
		"--format=%(authordate:unix) %(authorname) %(refname:short) %(HEAD)",
		"refs/heads",
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseGitBranches(string(output)), nil
}

func parseGitBranches(rawBranches string) *Branches {
	gb := newGitBranches()

	// parse git branch output by line
	sc := bufio.NewScanner(strings.NewReader(rawBranches))
	i := 0
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) < 3 {
			continue
		}
		b := Branch{
			Name:          fields[2],
			LastCommitter: fields[1],
			LastCommit:    parseUnixTimestamp(fields[0]),
		}

		if len(fields) == 4 && fields[3] == "*" {
			gb.CurrentBranch = b.Name
			gb.CurrentBranchIndex = i
		}

		gb.calcColumnWidth(b.LastCommitter, b.Name)
		gb.Branches = append(gb.Branches, b)

		i++
	}

	return gb
}

/*
	calculate the max width of each field to allow for even spacing
	and alignment across rows
*/
func (gb *Branches) calcColumnWidth(committer, name string) {
	if len(committer) > gb.MaxColumnWidth["committer"] {
		gb.MaxColumnWidth["committer"] = len(committer)
	}

	if len(name) > gb.MaxColumnWidth["name"] {
		gb.MaxColumnWidth["name"] = len(name)
	}
}

// FormatBranchStrings converts the branch obj to a pretty, formatted string
func (gb *Branches) FormatBranchStrings() []string {
	maxCommitterLen := gb.MaxColumnWidth["committer"]
	maxNameLen := gb.MaxColumnWidth["name"]

	items := make([]string, 0)
	for _, b := range gb.Branches {
		// TODO - add colors
		// TODO - add [current] tag
		formatted := fmt.Sprintf("%*s %*s %*s",
			-dateFieldWidth,
			b.formatDate(),
			-maxCommitterLen-committerFieldPadding,
			b.LastCommitter,
			-maxNameLen-nameFieldPadding,
			b.Name,
		)
		items = append(items, formatted)
	}
	return items
}
