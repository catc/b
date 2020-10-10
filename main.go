package main

import (
	"flag"
	"fmt"
	"os"
)

const usage = `By default, b checks out an existing branch. Usage:

Flags:
  -l         list all branches

Additional commands available:
  clone      create a clone of the current branch
  prune      select branches to remove
    -a       auto remove all branches older than 30 days (default: false)

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
	a := flag.Bool("a", false, "auto prune")
	l := flag.Bool("l", false, "list branches")
	flag.Parse()

	args := flag.Args()
	if len(args) > 0 {
		switch args[0] {
		case "prune":
			prune(*a)
			break
		case "clone":
			clone()
			break
		default:
			flag.Usage()
		}
	} else if *l {
		list()
	} else {
		checkout()
	}
}
