package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/catc/b/git"
	"github.com/eiannone/keyboard"
	"github.com/mgutz/ansi"
)

var green = ansi.ColorFunc("green+hb")
var red = ansi.ColorFunc("red+hb")
var white = ansi.ColorFunc("white+hb")
var yellowhb = ansi.ColorFunc("yellow")
var blue = ansi.ColorFunc("blue+hb")

const YES_KEY = "y"
const NO_KEY = "n"
const CHECKOUT_KEY = "c"
const DIFF_KEY = "d"

var SUPPORTED_KEYS = []string{YES_KEY, NO_KEY, CHECKOUT_KEY, DIFF_KEY}

func add() {
	files, err := git.GetChangedFiles()
	if err != nil {
		panic("Failed to get changed files")
	}
	if len(files) == 0 {
		fmt.Println("\nNo tracked files to diff")
	} else {
		fmt.Println(fmt.Sprintf("\n%v", green("Tracked files:")))
		for _, name := range files {
			judge(name, true)
		}
	}

	files, err = git.GetUntrackedFiles()
	if err != nil {
		panic("Failed to get untracked files")
	}
	if len(files) == 0 {
		fmt.Println("\nNo untracked files to diff")
	} else {
		fmt.Println(fmt.Sprintf("\n%v", red("Untracked files:")))
		for _, name := range files {
			judge(name, false)
		}
		fmt.Println()
	}
}

func judge(name string, tracked bool) {
	view := func() {
		if err := viewDiff(name, tracked); err != nil {
			fmt.Println("Failed to view file", err)
			os.Exit(0)
			return
		}
	}

	// show diff
	view()

	// post question
	fmt.Printf(`%v %v   %v `,
		white("Add:"),
		yellowhb(name),
		blue(`[y\n\d\c]`),
	)

	status := "noop"
	// wait for user input
	for {
		key := waitForKey()
		if key == DIFF_KEY {
			view()
		} else {
			if key == CHECKOUT_KEY {
				status = "checkout"
			}
			if key == YES_KEY {
				status = "add"
			}
			break
		}
	}

	// add file
	if status == "add" {
		cmd := exec.Command("git", "add", name)
		if err := cmd.Run(); err != nil {
			fmt.Println(fmt.Sprintf(`Failed to add file "%v"`, name))
			panic(err)
		}
	}

	if status == "checkout" {
		checkoutFile(name)
	}

	// update output
	var prefix string
	if status == "checkout" {
		prefix = "🛒"
	} else if status == "add" {
		prefix = "✅"
	} else {
		prefix = "❌"
	}
	fmt.Printf("\r%v   %v%v\n",
		prefix,
		yellowhb(name),
		// for padding to clear line
		strings.Repeat(" ", 12),
	)
}

func viewDiff(name string, tracked bool) error {
	var cmd *exec.Cmd

	if tracked {
		// `git diff HEAD -- FILEPATH` - handles removed files as well
		cmd = exec.Command("git", "diff", "--color=always", "HEAD", "--", name)
	} else {
		cmd = exec.Command("git", "diff", "--color=always", "--no-index", "/dev/null", name)
	}

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	// for untracked files, ignore error - git will exit with status 1 since no diffs
	if err := cmd.Run(); err != nil && tracked {
		fmt.Println(fmt.Sprintf(`Failed to view file "%v"`, name))
		return err
	}
	return nil
}

func waitForKey() string {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}

	close := func() {
		if err := keyboard.Close(); err != nil {
			panic(err)
		}
	}
	defer close()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		// handle exit
		if key == keyboard.KeyCtrlC || key == keyboard.KeyCtrlD {
			close()
			os.Exit(130)
		}

		// handle y/n/d/c
		c := strings.ToLower(string(char))
		for _, v := range SUPPORTED_KEYS {
			if v == c {
				return c
			}
		}
	}
}

func checkoutFile(filepath string) {
	cmd := exec.Command("git", "checkout", filepath)
	if err := cmd.Run(); err != nil {
		fmt.Println(fmt.Sprintf(`Failed to checkout file "%v"`, filepath))
		panic(err)
	}
}
