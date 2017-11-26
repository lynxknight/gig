package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/pkg/term"
)

var highlighter = color.New(color.BgWhite, color.FgBlack).SprintfFunc()

func drawUI(branches []branch, query string, cursorpos int) {
	clearScreen()
	fmt.Println(query)
	fmt.Println("============")
	displayBranches(branches, cursorpos)
}

func displayBranches(branches []branch, cursorpos int) {
	var name string
	for index, branch := range branches {
		if cursorpos == index {
			name = highlighter("%v", branch.name)
		} else {
			name = branch.name
		}
		fmt.Println(name)
	}
}

type MGMT int

const ( // Do not handle a lot of stuff since there is no cursor concept
	MGMT_TEXT MGMT = iota

	MGMT_CR

	MGMT_ARROW_UP
	MGMT_ARROW_DOWN

	MGMT_CTRL_C
	MGMT_CTRL_D
	MGMT_CTRL_W

	MGMT_BACKSPACE
)

type termInput struct {
	rawValue []byte
	mgmt     MGMT
}

// getTermInput normally returns parsed termInput, upon unknown stuff returns
// error
func getTermInput() (result termInput, err error) {
	numRead, bytes, err := readTerm()
	if err != nil {
		return result, err
	}
	result.rawValue = bytes[:numRead]
	if numRead == 3 && bytes[0] == 27 && bytes[1] == 91 {
		// Three-character control sequence, beginning with "ESC-[".
		// Since there are no ASCII codes for arrow keys, we use
		// Javascript key codes.
		if bytes[2] == 65 {
			result.mgmt = MGMT_ARROW_UP
		} else if bytes[2] == 66 {
			result.mgmt = MGMT_ARROW_DOWN
		}
	} else if numRead == 1 {
		ascii := int(bytes[0])
		switch ascii {
		case 3:
			result.mgmt = MGMT_CTRL_C
		case 4:
			result.mgmt = MGMT_CTRL_D
			os.Exit(0)
		case 13:
			result.mgmt = MGMT_CR
		case 23:
			result.mgmt = MGMT_CTRL_W
		case 127:
			result.mgmt = MGMT_BACKSPACE
		}
	}
	return result, nil
}

func readTerm() (numRead int, bytes []byte, err error) {
	t, _ := term.Open("/dev/tty")
	term.RawMode(t)
	bytes = make([]byte, 140)
	numRead, err = t.Read(bytes)
	if err != nil {
		return
	}
	t.Restore()
	t.Close()
	return
}

func clearScreen() {
	if os.Getenv("NOCLEAR") == "1" {
		return
	}
	print("\033[H\033[2J")
}
