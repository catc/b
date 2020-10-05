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
