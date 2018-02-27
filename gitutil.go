package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type ref struct {
	name  string
	mdate time.Time
}

type costcacheT map[string]int

type branch struct {
	name      string
	costcache costcacheT
}

func isGitRoot(dir string) bool {
	direntries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln("Failed to read directory", dir, err)
	}
	for _, d := range direntries {
		if d.Name() == ".git" {
			return true
		}
		// ReadDir is sorted, "." should be first, if not => we lost
		if !strings.HasPrefix(d.Name(), ".") {
			return false
		}
	}
	return false
}

func getGitRoot() string {
	// get cwd
	// move up until you notice .git folder
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("Failed to get cwd:", err)
	}
	for !isGitRoot(cwd) && cwd != "/" {
		cwd = path.Dir(cwd)
	}
	if cwd == "/" && !isGitRoot(cwd) {
		log.Fatalln("fatal: not a git repository (or any of the parent directories): .git")
	}
	return path.Join(cwd, ".git")
}

func getRefsRec(prefix, dir string) []ref {
	s := make([]ref, 0)
	direntries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln("Failed to read directory", dir, err)
	}
	newPrefix := path.Join(prefix, dir)
	for _, d := range direntries {
		if d.IsDir() {
			s = append(
				s, getRefsRec(newPrefix, path.Join(newPrefix, d.Name()))...,
			)
			continue
		}
		s = append(
			s, ref{name: path.Join(newPrefix, d.Name()), mdate: d.ModTime()},
		)
	}
	return s
}

func getRefs(gitRoot string) []ref {
	refsPath := path.Join(gitRoot, "refs", "heads")
	return getRefsRec("", refsPath)
}

func getBranches() []branch {
	// git for-each-ref is unstable across git versions, so we implement it
	gitRoot := getGitRoot()
	refs := make([]ref, 0)
	refsPath := path.Join(gitRoot, "refs", "heads")
	// TODO: refs/ provides full path, so we could not checkout into this
	// refs := getRefs(gitRoot)
	filepath.Walk(refsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalln("Runtime failed to walk on path", path, err)
		}
		if !info.IsDir() {
			refs = append(refs, ref{name: info.Name(), mdate: info.ModTime()})
		}
		return nil
	})
	branches := make([]branch, len(refs))
	for i := range refs {
		branches[i] = branch{
			name:      refs[i].name,
			costcache: make(costcacheT),
		}
		// Empty QS = sort by date
		// Screw you if you run on 32bit system
		branches[i].costcache[""] = -int(refs[i].mdate.Unix())
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
