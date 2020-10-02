package git

import (
	"bufio"
	"fmt"
	"os/exec"
	"reflect"
	"strings"
)

// Branches contains all branches and git related configs
type Branches struct {
	Branches           []Branch
	CurrentBranch      string
	CurrentBranchIndex int
	MaxColumnLength    map[string]int
}

// GetBranches fetches all branches
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

	gb := Branches{
		Branches:        make([]Branch, 0),
		MaxColumnLength: make(map[string]int),
	}

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

		gb.Branches = append(gb.Branches, b)
	}

	return &gb
}

// TODO
func (branches *Branches) getBranchMetadata() {
	for _, b := range branches.Branches {
		r := reflect.ValueOf(b)
		for _, field := range fieldsToDisplay {
			fmt.Println(field, r)
		}
	}
}

// FormatBranchStrings returns the branch strings formatted for survey
func (branches *Branches) FormatBranchStrings() []string {
	items := make([]string, 0)
	for _, b := range branches.Branches {
		formatted := fmt.Sprintf("%-11v %v %v", b.formatDate(), b.LastCommitter, b.Name)
		items = append(items, formatted)
	}
	return items
}
