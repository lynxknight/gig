package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type branch struct {
	name      string
	costcache map[string]int
}

func getBranches() []branch {
	gbOutput, err := exec.Command(
		"git", "branch",
		"--sort", "-committerdate",
		"--format", "%(refname:short)|%(creatordate:unix)",
	).Output()
	if err != nil {
		panic(err)
	}
	rawBranches := strings.Split(string(gbOutput), "\n")
	// Last line contains emptyspace, strip it via -1 in make() and for{}
	branches := make([]branch, len(rawBranches)-1)
	for i := 0; i < len(rawBranches)-1; i++ {
		splitted := strings.Split(rawBranches[i], "|")
		timestamp, _ := strconv.Atoi(splitted[1])
		costcache := make(map[string]int)
		costcache[""] = timestamp // Empty querystring = sort by date
		branches[i] = branch{
			name:      splitted[0],
			costcache: costcache,
		}
	}
	return branches
}

func checkoutBranch(branch string) error {
	gcheckout := exec.Command("git", "checkout", branch)
	out, err := gcheckout.Output()
	if err != nil {
		return fmt.Errorf(string(err.(*exec.ExitError).Stderr))
	}
	fmt.Println(string(out))
	return nil
}
