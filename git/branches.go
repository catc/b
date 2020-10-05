package git

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mgutz/ansi"
)

const dateFieldWidth = 12
const committerFieldPadding = 5

// +4 to account for asterisk padding on current branch
const nameFieldPadding = 1 + 4

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
		1601742631 author_name master *
	*/
	cmd := exec.Command(
		gitExecutable,
		"for-each-ref",
		"--sort=-authordate:iso8601",
		"--format=%(authordate:unix)|%(authorname)|%(refname:short)|%(HEAD)",
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
		fields := strings.Split(sc.Text(), "|")
		if len(fields) < 4 {
			continue
		}
		b := Branch{
			Name:          fields[2],
			LastCommitter: fields[1],
			LastCommit:    parseUnixTimestamp(fields[0]),
		}

		if fields[3] == "*" {
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

	// colors
	dateFormat := ansi.ColorFunc("green+hb")
	committerFormat := ansi.ColorFunc("white")
	nameFormat := ansi.ColorFunc("white+hb")
	asteriskFormat := ansi.ColorFunc("5+hb")

	items := make([]string, 0)
	for _, b := range gb.Branches {
		name := nameFormat(b.Name)
		// pad current branch with asterisks
		if b.Name == gb.CurrentBranch {
			name = asteriskFormat("✱ ") +
				name +
				asteriskFormat(" ✱")
		}

		// column alignment
		name = fmt.Sprintf("%*s ", -maxNameLen-nameFieldPadding, name)
		date := fmt.Sprintf("%*s", -dateFieldWidth, b.formatDate())
		committer := fmt.Sprintf("%*s", -maxCommitterLen-committerFieldPadding, b.LastCommitter)

		formatted := dateFormat(date) +
			committerFormat(committer) +
			name

		items = append(items, formatted)
	}

	return items
}

// ChangeBranch switches branch
func ChangeBranch(branch string) {
	gitExecutable, _ := exec.LookPath("git")
	cmdGoVer := &exec.Cmd{
		Path:   gitExecutable,
		Args:   []string{gitExecutable, "checkout", branch},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}
	cmdGoVer.Run()
}
