package main

import (
	"fmt"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/catc/b/git"
	"github.com/mgutz/ansi"
)

const day = time.Hour * 24

func prune(days int) {
	gb, err := git.GetBranches()
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(gb.Branches) == 1 {
		fmt.Println("This repo has no branches")
		return
	}

	if days > 0 {
		autoPrune(gb, days)
		return
	}

	displayIntro(gb)

	selected := []int{}
	options := gb.FormatBranchStrings(false)
	prompt := &survey.MultiSelect{
		Message:  "Select branches to delete:",
		Options:  options,
		PageSize: pageSize,
	}
	survey.AskOne(prompt, &selected)

	confirm(gb, selected)
}

func confirm(gb *git.Branches, selected []int) {
	if len(selected) == 0 {
		fmt.Println("No branches selected to delete")
		return
	}

	// display branches to be deleted
	options := gb.FormatBranchStrings(true)
	answers := "\n" + ansi.Color("Selected:\n", "yellow+hb")
	for _, i := range selected {
		answers += fmt.Sprintf("%s\n", options[i])
	}
	fmt.Println(answers)

	// confirm deletion
	delete := false
	prompt := &survey.Confirm{
		Message: "Are you sure you want to delete these branches?",
	}
	survey.AskOne(prompt, &delete)
	if !delete {
		return
	}

	// get branch names
	branches := []string{}
	for _, i := range selected {
		branch := gb.Branches[i].Name
		branches = append(branches, branch)
	}

	git.DeleteBranches(branches)
}

func autoPrune(gb *git.Branches, days int) {
	cutoff := time.Now().Add(-day * time.Duration(days))

	selected := []int{}
	for i, b := range gb.Branches {
		if b.LastCommit.Before(cutoff) {
			selected = append(selected, i)
		}
	}

	confirm(gb, selected)
}
