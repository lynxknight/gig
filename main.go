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
	if len(branches) == 0 {
		fmt.Println("No branches found.")
		return
	}
	if !exactMatch(target, branches) {
		target = pickBranch(target, branches).name
	}
	if err := checkoutBranch(target); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func pickBranch(target string, branches []branch) branch {
	terminal := getTerm() // TODO: global variable in UI?

	cursorPos := 0
	var queryStringBuf bytes.Buffer
	queryStringBuf.WriteString(target)

	sortBranches(branches, queryStringBuf.String())
	drawUI(branches, queryStringBuf.String(), cursorPos, getTermHeight())
	for {
		resort := true
		moveCursor(1, queryStringBuf.Len()+1)
		usrInput, _ := getUserInput(terminal)
		switch usrInput.input {
		case inputArrowUp:
			if cursorPos > 0 {
				cursorPos--
			}
			resort = false
		case inputArrowDown:
			if cursorPos < len(branches)-1 {
				cursorPos++
			}
			resort = false // TODO: Redraw only what's needed
		case inputCtrlW:
			queryStringBuf.Truncate(0)
		case inputCtrlC:
			clearScreen()
			os.Exit(1)
		case inputCtrlD:
			clearScreen()
			os.Exit(0)
		case inputBackspace:
			if queryStringBuf.Len() > 0 {
				queryStringBuf.Truncate(queryStringBuf.Len() - 1)
			} else {
				continue
			}
		case inputText:
			cursorPos = 0
			queryStringBuf.Write(usrInput.rawValue) // TODO: handle errors?
		case inputCR:
			clearScreen()
			return branches[cursorPos]
		default:
			continue
		}
		if resort {
			sortBranches(branches, queryStringBuf.String())
		}
		drawUI(branches, queryStringBuf.String(), cursorPos, getTermHeight())
	}
}

func sortBranches(branches []branch, query string) {
	// Calculate distance for querystring if we have not done it yet
	// start := time.Now().UnixNano()
	if _, ok := branches[0].costcache[query]; !ok {
		for i := range branches {
			branches[i].costcache[query] = distance.GetScore(branches[i].name, query)
		}
	}
	// start2 := time.Now().UnixNano()
	// cost := start2 - start
	sort.SliceStable(branches, func(i, j int) bool {
		return branches[i].costcache[query] <= branches[j].costcache[query]
	})
	// sort := time.Now().UnixNano() - start2
	// total := sort + cost
	moveCursor(2, 0)
	// f, err := os.OpenFile("/tmp/metrix.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	// if err != nil {
	// log.Fatalln(err)
	// }
	// defer f.Close()
	// w, err := f.WriteString(fmt.Sprintf("cost: %v; sort: %v; total: %v\n", cost, sort, total))
	// fmt.Println(w)
	// if err != nil {
	// log.Fatalln(err)
	// }
}
