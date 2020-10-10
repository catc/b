package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/catc/b/git"
	"github.com/mgutz/ansi"
)

const pageSize = 18

// checkout existing branch
func checkout() {
	gb, err := git.GetBranches()
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(gb.Branches) == 1 {
		fmt.Println("This repo has no branches")
		return
	}

	displayIntro(gb)

	prompt := &survey.Select{
		Message: ansi.Color(" Select a branch:", "white+b"),
		Options: gb.FormatBranchStrings(false),
		Filter: func(filter string, value string, i int) bool {
			name := strings.ToLower(gb.Branches[i].Name)
			return strings.Contains(name, strings.ToLower(filter))
		},
		PageSize: pageSize,
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
		survey.WithIcons(setIcons),
	); err == terminal.InterruptErr {
		os.Exit(0)
	} else if err != nil {
		fmt.Println("Error with survey:", err)
		return
	}

	git.ChangeBranch(gb.Branches[index].Name)
}

func setIcons(icons *survey.IconSet) {
	icons.SelectFocus.Text = "🌿"
}
