package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/catc/b/git"
)

func changeBranch() {
	gb, err := git.GetBranches()
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(gb.Branches) == 1 {
		fmt.Println("This repo has no branches")
		return
	}

	prompt := &survey.Select{
		Message: "Select a branch",
		Options: gb.FormatBranchStrings(),
		Filter: func(filter string, value string, i int) bool {
			name := gb.Branches[i].Name
			return strings.Contains(name, filter)
		},
		PageSize: 10,
	}

	differentBranchValidator := func(val interface{}) error {
		option, ok := val.(survey.OptionAnswer)
		if !ok || option.Index == gb.CurrentBranchIndex {
			return errors.New("You are currently on this branch")
		}
		return nil
	}

	var index int
	if err = survey.AskOne(
		prompt,
		&index,
		survey.WithValidator(differentBranchValidator),
	); err == terminal.InterruptErr {
		os.Exit(0)
	} else if err != nil {
		fmt.Println("Error with survey", err)
		return
	}

	fmt.Println(index)

	// TODO - change to selected branch and exit
}
