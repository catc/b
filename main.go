package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

const usage = `By default, b checks out an existing branch. Usage:

Flags:
  -l          list all branches

Additional commands available:
  clone       create a clone of the current branch
	prune [#]   select branches to remove, can optionally specify how many days in past to prune
  add         interactively stage files
`

func init() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	// permutate args (so `b prune -a` works correctly)
	// https://stackoverflow.com/questions/25113313/golang-flag-stops-parsing-after-the-first-non-option
	args := os.Args[1:]
	optind := 0
	for i := range args {
		if args[i][0] == '-' {
			tmp := args[i]
			args[i] = args[optind]
			args[optind] = tmp
			optind++
		}
	}
}

func main() {
	lFlag := flag.Bool("l", false, "list branches")
	flag.Parse()

	args := flag.Args()
	if len(args) > 0 {
		switch args[0] {
		case "prune":
			// prune by date (ie: days in past)
			if len(args) == 2 {
				i, err := strconv.Atoi(args[1])
				if err != nil {
					flag.Usage()
					return
				}
				prune(i)
				return
			}
			// manually select branches to prune
			prune(0)
			break
		case "clone":
			clone()
			break
		case "add":
			add()
			break
		default:
			flag.Usage()
		}
	} else if *lFlag {
		list()
	} else {
		checkout()
	}
}
