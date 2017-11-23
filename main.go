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
		ascii, keyCode, _ := getChar()
		if keyCode != 0 {
			switch keyCode {
			case 38: // Up
				if cursorpos > 0 {
					cursorpos--
				}
			case 40: // Down
				if cursorpos < len(branches)-1 {
					cursorpos++
				}
			}
		} else {
			switch ascii {
			case 13: // <CR>
				clearScreen()
				return branches[cursorpos]
			case 3: // ctrl-c
				clearScreen()
				os.Exit(1)
			case 4: // ctrl-d
				os.Exit(0)
			case 127: // backspace
				if buffer.Len() > 0 {
					buffer.Truncate(buffer.Len() - 1)
				}
			case 23: // ctrl-w
				buffer.Truncate(0)
			default:
				char := string(ascii)
				buffer.WriteString(char)
			}
		}
		sortBranches(branches, buffer.String())
		drawUI(branches, buffer.String(), cursorpos)
	}
}

func sortBranches(branches []branch, query string) {
	// Calculate distance for querystring if we have not done it yet
	if _, ok := branches[0].costcache[query]; !ok {
		for i := range branches {
			branches[i].costcache[query] = distance.LevenshteinDistance(query, branches[i].name)
		}
	}
	sort.Slice(branches, func(i, j int) bool {
		return branches[i].costcache[query] < branches[j].costcache[query]
	})

}
