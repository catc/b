package main

import (
	"fmt"

	"github.com/catc/b/git"
	"github.com/mgutz/ansi"
)

func displayIntro(gb *git.Branches) {
	intro := fmt.Sprintf("\n%v local branches:\n", len(gb.Branches))
	fmt.Println(ansi.Color(intro, "white+bh"))
}
