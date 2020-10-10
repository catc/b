package main

import (
	"fmt"

	"github.com/catc/b/git"
)

func list() {
	gb, err := git.GetBranches()
	if err != nil {
		fmt.Println(err)
		return
	}

	displayIntro(gb)

	for _, b := range gb.FormatBranchStrings(true) {
		fmt.Println(b)
	}
	fmt.Println("")
}
