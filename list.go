package main

import (
	"fmt"

	"github.com/catc/b/git"
	"github.com/mgutz/ansi"
)

func list() {
	gb, err := git.GetBranches()
	if err != nil {
		fmt.Println(err)
		return
	}

	intro := fmt.Sprintf("%v local branches:", len(gb.Branches))
	fmt.Println(ansi.Color(intro, "white+bh"))

	for _, b := range gb.FormatBranchStrings() {
		fmt.Println(b)
	}
}
