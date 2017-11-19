package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/lynxknight/gig/distance"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/pkg/term"
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
	branches = branches[:len(branches)-1] // Strip last space
	return branches
}

func checkoutBranch(branch string) {
	gcheckout := exec.Command("git", "checkout", branch) // TO AFRAID TO CHECKOUT
	out, err := gcheckout.Output()
	fmt.Println(string(out))
	if err != nil {
		fmt.Println(string(err.(*exec.ExitError).Stderr))
		log.Fatalln("Failed to checkout branch, reason:", err)
	}
}

func exactMatch(target string, branches []string) bool {
	for _, branch := range branches {
		if branch == target {
			return true
		}
	}
	return false
}

func showLev(target string, branches []string, highlight int) string {
	choice := ""
	type container struct {
		name string
		cost int
	}
	icb := make([]container, len(branches))
	for i, branch := range branches {
		icb[i] = container{branch, distance.Distance(target, branch)}
	}

	if target != "" {
		sort.Slice(icb, func(i, j int) bool {
			return icb[i].cost < icb[j].cost
		})
	}
	if highlight > -1 && highlight < len(branches) {
		choice = icb[highlight].name
		icb[highlight].name = color.New(color.BgWhite, color.FgBlack).SprintfFunc()("%v", icb[highlight].name)
	}
	for _, cont := range icb {
		fmt.Println(cont.name)
	}
	return choice
}

func main() {
	assureStdoutIsTTY()
	target := parseArgs()
	branches := getBranches()
	if !exactMatch(target, branches) {
		target = pickBranch(target, branches)
	}
	checkoutBranch(target)
}

func pickBranch(target string, branches []string) string {
	var choice string
	var buffer bytes.Buffer
	buffer.WriteString(target)
	cursorpos := 0
	current := target
	choice = showLev(current, branches, cursorpos)
	// Enter REPL
	for {
		ascii, keyCode, _ := getChar()
		fmt.Println(ascii, keyCode)
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
				fmt.Println("RETURN")
				return choice
			case 3: // ctrl-c
				fmt.Println("INTERRUPT")
				os.Exit(1)
			case 4: // ctrl-d
				fmt.Println("EOF")
				os.Exit(0)
			case 127: // backspace
				fmt.Println("EOF")
				if buffer.Len() > 0 {
					buffer.Truncate(buffer.Len() - 1)
				}
			case 23: // ctrl-w
				buffer.Truncate(0)
			default:
				char := string(ascii)
				buffer.WriteString(char)
			}
			fmt.Println(ascii)
		}

		clearScreen()
		current = buffer.String()
		fmt.Println(current)
		choice = showLev(current, branches, cursorpos)
	}
}
func clearScreen() {
	print("\033[H\033[2J")
}

func getChar() (ascii int, keyCode int, err error) {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes := make([]byte, 3)

	var numRead int
	numRead, err = t.Read(bytes)
	if err != nil {
		return
	}
	if numRead == 3 && bytes[0] == 27 && bytes[1] == 91 {
		// Three-character control sequence, beginning with "ESC-[".

		// Since there are no ASCII codes for arrow keys, we use
		// Javascript key codes.
		if bytes[2] == 65 {
			// Up
			keyCode = 38
		} else if bytes[2] == 66 {
			// Down
			keyCode = 40
		} else if bytes[2] == 67 {
			// Right
			keyCode = 39
		} else if bytes[2] == 68 {
			// Left
			keyCode = 37
		}
	} else if numRead == 1 {
		ascii = int(bytes[0])
	} else {
		// Two characters read??
	}
	t.Restore()
	t.Close()
	return
}
