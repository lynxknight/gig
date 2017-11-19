package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

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

func getBranches() []string {
	gbOutput, err := exec.Command(
		"git", "branch",
		"--sort", "-committerdate",
		"--format", "%(refname:short)",
	).Output()
	if err != nil {
		panic(err)
	}
	branches := strings.Split(string(gbOutput), "\n")
	return branches
}

func exactMatch(target string, branches []string) bool {
	for _, branch := range branches {
		if branch == target {
			return true
		}
	}
	return false
}

func showLev(target string, branches []string) {
	type container struct {
		name string
		cost int
	}
	icb := make([]container, len(branches))
	for i, branch := range branches {
		icb[i] = container{branch, distance.Distance(target, branch)}
	}
	sort.Slice(icb, func(i, j int) bool {
		return icb[i].cost < icb[j].cost
	})
	for _, cont := range icb {
		fmt.Println(cont.name, cont.cost)
	}
}

func main() {
	assureStdoutIsTTY()
	target := parseArgs()
	branches := getBranches()
	if exactMatch(target, branches) {
		out, err := exec.Command("git", "checkout", target).Output()
		if err != nil {
			log.Fatalln("Failed to checkout branch")
		}
		fmt.Print(out)
		return
	}
	showLev(target, branches)
}
