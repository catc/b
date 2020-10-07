package main

import (
	"fmt"
	"os"

	"github.com/catc/b/git"
	"github.com/mgutz/ansi"
)

func clone() {
	gb, err := git.GetBranches()
	if err != nil {
		fmt.Println(err)
		return
	}

	if gb.CurrentBranch == "" {
		fmt.Fprint(os.Stderr, "Failed to detect current branch")
		return
	}

	branches := make(map[string]bool)
	for _, b := range gb.Branches {
		branches[b.Name] = true
	}

	// attempt to find a name for the cloned branch
	for i := 1; i < 10; i++ {
		// get branch clone name
		clone := getCloneName(gb.CurrentBranch, i)

		// create the branch if it doesn't exist with current name
		if !branches[clone] {
			git.CreateBranch(clone)
			msg := fmt.Sprintf("Creating branch \"%s\"", ansi.Color(clone, "green+hb"))
			fmt.Println(msg)
			break
		}
	}
}

func getCloneName(branch string, index int) string {
	return fmt.Sprintf("%v__CLONE-%v", branch, index)
}
