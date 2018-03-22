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

func getPackedRefs(gitRoot string) ([]ref, error) {
	// git has a "more perfomant" way of storing refs info, called "packed-refs"
	// more info at https://git-scm.com/docs/git-pack-refs
	byteContent, err := ioutil.ReadFile(path.Join(gitRoot, "packed-refs"))
	if err != nil {
		return nil, err
	}
	content := string(byteContent)
	refsHeads := "refs/heads" // we care about heads only now
	refs := make([]ref, 0)
	for _, line := range strings.Split(content, "\n") {
		objects := strings.Split(line, " ")
		if len(objects) != 2 { // if it is a ^SHA record...
			continue // it is not interesting to us
		}
		refName := objects[1]
		if strings.HasPrefix(refName, refsHeads) {
			refs = append(
				refs,
				ref{name: refName[len(refsHeads)+1:], mdate: time.Unix(0, 0)},
			)
		}
	}
	return refs, nil
}

func getBranches() []branch {
	// git for-each-ref is unstable across git versions, so we implement it
	gitRoot := getGitRoot()
	refs := make([]ref, 0)
	refsPath := path.Join(gitRoot, "refs", "heads")
	filepath.Walk(refsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalln("Runtime failed to walk on path", path, err)
		}
		if !info.IsDir() {
			// TODO: what if ref is presented in both refs and packed-refs?
			r := ref{name: path[len(refsPath)+1:], mdate: info.ModTime()}
			refs = append(refs, r)
		}
		return nil
	})
	packedRefs, err := getPackedRefs(gitRoot)
	if err == nil {
		refs = append(refs, packedRefs...)
	}
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
