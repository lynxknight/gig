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
	log.Println("Got", len(branches), "branches")
	return branches
}

func exactMatch(target string, branches []string) bool {
	for _, branch := range branches {
		log.Println("Checkin if", branch, "is equal to", target)
		if branch == target {
			log.Println("Exact match on", target)
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
	branches := getBranches()
	target := os.Args[1]
	log.Println("Target:", target)
	if exactMatch(target, branches) {
		log.Println("Checking out", target)
		_, err := exec.Command("git", "checkout", target).Output()
		if err != nil {
			log.Fatalln("Failed to checkout branch")
		}
		return
	}
	showLev(target, branches)
}
