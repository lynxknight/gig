package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/lynxknight/gig/distance"

	"github.com/mattn/go-isatty"
)

func assureStdoutIsTTY() {
	stdin, stdout := os.Stdin.Fd(), os.Stdout.Fd()
	if !(isatty.IsTerminal(stdin) && isatty.IsTerminal(stdout)) {
		log.Fatalln("Supposed to be run using TTY")
	}
}

func parseArgs() string {
	argc := len(os.Args)
	if argc > 2 {
		fmt.Println("Usage:", os.Args[0], "[branch]")
		os.Exit(1)
	}
	if argc == 2 {
		return os.Args[1]
	}
	return ""
}

func exactMatch(target string, branches []branch) bool {
	for _, branch := range branches {
		if branch.name == target {
			return true
		}
	}
	return false
}

func main() {
	assureStdoutIsTTY()
	target := parseArgs()
	branches := getBranches()
	if !exactMatch(target, branches) {
		target = pickBranch(target, branches).name
	}
	if err := checkoutBranch(target); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func pickBranch(target string, branches []branch) branch {
	var buffer bytes.Buffer // buffer contains querystring
	buffer.WriteString(target)
	cursorpos := 0 // cursorpos stores current menu position
	sortBranches(branches, buffer.String())
	drawUI(branches, buffer.String(), cursorpos)
	for { // REPL
		input, _ := getTermInput()
		switch input.mgmt {
		case MGMT_ARROW_UP: // Up
			if cursorpos > 0 {
				cursorpos--
			}
		case MGMT_ARROW_DOWN: // Up
			if cursorpos < len(branches)-1 {
				cursorpos++
			}
		case MGMT_CTRL_W:
			buffer.Truncate(0)
		case MGMT_CTRL_C:
			clearScreen()
			os.Exit(1)
		case MGMT_CTRL_D:
			clearScreen()
			os.Exit(0)
		case MGMT_BACKSPACE:
			if buffer.Len() > 0 {
				buffer.Truncate(buffer.Len() - 1)
			}
		case MGMT_TEXT:
			buffer.Write(input.rawValue) // TODO: handle errors?
		}
		// TODO: sometimes we don't need to resort branches
		sortBranches(branches, buffer.String())
		drawUI(branches, buffer.String(), cursorpos)
	}
}

func sortBranches(branches []branch, query string) {
	// Calculate distance for querystring if we have not done it yet
	if _, ok := branches[0].costcache[query]; !ok {
		for i := range branches {
			branches[i].costcache[query] = distance.GetScore(branches[i].name, query)
		}
	}
	sort.Slice(branches, func(i, j int) bool {
		return branches[i].costcache[query] < branches[j].costcache[query]
	})

}
