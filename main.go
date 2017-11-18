package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

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
	}
}
