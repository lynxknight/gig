package main

import (
	"os/exec"
)

func main() {
	BranchCmd := exec.Command("git", "branch")
	branch, err := BranchCmd.Output()
	if err != nil {
		panic(err)
	}
	print(string(branch))
}
