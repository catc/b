package main

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/catc/b/git"
)

func main() {
	// TODO - add flags + options
	changeBranch()
}

func changeBranch() {
	b, err := git.GetBranches()
	if err != nil {
		fmt.Println("Error fetching branches -", err.Error())
		return
	}

	if len(b.Branches) == 1 {
		fmt.Println("This repo has no branches")
		return
	}

	prompt := &survey.Select{
		Message: "Select a branch:",
		// Options:  []string{"\033[0;33mred\033[0m and \033[0;32mremaining", "blue", "green"},
		Options: b.FormatBranchStrings(),
		Filter: func(filter string, value string, i int) bool {
			name := b.Branches[i].Name
			return strings.Contains(name, filter)
		},
		PageSize: 10,
	}

	var index int
	if err = survey.AskOne(prompt, &index); err != nil {
		fmt.Println("Error with survey", err)
		return
	}

	fmt.Println(index)

}
