package git

import (
	"bufio"
	"fmt"
	"os/exec"
	"sort"
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

	cmd := exec.Command(
		gitExecutable,
		"--no-pager", "branch",
	)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return parseGitBranches(string(output)), nil
}

func parseGitBranches(rawBranches string) *Branches {
	// parse git branch output by line
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(rawBranches))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	gb := newGitBranches()

	// get individual branch name + metadata
	for i, str := range lines {
		fields := strings.Fields(str)

		var name string
		if len(fields) > 0 && fields[0] == "*" {
			name = fields[1]
			gb.CurrentBranch = name
			gb.CurrentBranchIndex = i
		} else {
			name = fields[0]
		}

		b := Branch{Name: name}
		b.populateBranchMetadata()

		gb.calcColumnWidth(b.LastCommitter, b.Name)
		gb.Branches = append(gb.Branches, b)
	}

	// sort by date
	branches := gb.Branches
	sort.Slice(gb.Branches, func(i, j int) bool {
		return branches[i].LastCommit.After(branches[j].LastCommit)
	})

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
