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

type inputType int

const ( // Do not handle a lot of stuff since there is no cursor concept
	inputText inputType = iota

	inputCR

	inputArrowUp
	inputArrowDown

	inputCtrlC
	inputCtrlD
	inputCtrlW

	inputBackspace
)

type userInput struct {
	rawValue []byte
	input     inputType
}

// getUserInput normally returns parsed userInput, upon unknown stuff returns
// error
func getUserInput() (result userInput, err error) {
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
			result.input = inputArrowUp
		} else if bytes[2] == 66 {
			result.input = inputArrowDown
		}
	} else if numRead == 1 {
		ascii := int(bytes[0])
		switch ascii {
		case 3:
			result.input = inputCtrlC
		case 4:
			result.input = inputCtrlD
			os.Exit(0)
		case 13:
			result.input = inputCR
		case 23:
			result.input = inputCtrlW
		case 127:
			result.input = inputBackspace
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
