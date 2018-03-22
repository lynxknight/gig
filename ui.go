package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/pkg/term"
	"golang.org/x/crypto/ssh/terminal"
)

func min(ints ...int) int {
	m := ints[0]
	for i := 1; i < len(ints); i++ {
		if ints[i] < m {
			m = ints[i]
		}
	}
	return m
}

var highlighter = color.New(color.BgWhite, color.FgBlack).SprintfFunc()

func drawUI(branches []branch, query string, cursorpos, height int) {
	clearScreen()
	fmt.Print(query + "\r\n")
	fmt.Print("============" + "\r\n")
	displayBranches(branches, cursorpos, height)
}

func displayBranches(branches []branch, cursorpos, height int) {
	// var name string
	branchesToPrint := make([]string, min(height-2, len(branches)))
	for i := range branchesToPrint {
		if cursorpos == i {
			branchesToPrint[i] = highlighter("%v", branches[i].name)
		} else {
			branchesToPrint[i] = branches[i].name
		}
	}
	fmt.Print(strings.Join(branchesToPrint, "\r\n"))
	// for index, branch := range branches {
	// if index+4 > height {
	// break
	// }
	// if cursorpos == index {
	// name = highlighter("%v", branch.name)
	// } else {
	// name = branch.name
	// }
	// fmt.Println(name)
	// }
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

	inputOther
)

type userInput struct {
	rawValue []byte
	input    inputType
}

// getUserInput normally returns parsed userInput, upon unknown stuff returns
// inputMeta (which is not Meta key, but "something I don't know what to do with")
func getUserInput(terminal *term.Term) (result userInput, err error) {
	numRead, bytes, err := readTerm(terminal)
	if err != nil {
		return result, err
	}
	result.rawValue = bytes[:numRead]
	if numRead == 3 && bytes[0] == 27 && bytes[1] == 91 {
		if bytes[2] == 65 {
			result.input = inputArrowUp
		} else if bytes[2] == 66 {
			result.input = inputArrowDown
		} else {
			result.input = inputOther
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
		default:
			if ascii < 32 || ascii > 126 {
				result.input = inputOther
			}
		}
	} else {
		result.input = inputOther
	}
	return result, nil
}

func readTerm(t *term.Term) (numRead int, bytes []byte, err error) {
	bytes = make([]byte, 140)
	numRead, err = t.Read(bytes)
	if err != nil {
		return
	}
	return
}

func getTerm() *term.Term {
	t, err := term.Open("/dev/tty")
	term.RawMode(t)
	if err != nil {
		panic("Failed to open tty device")
	}
	print("\033[?1049h")
	return t
}

func restoreTerm(terminal *term.Term) {
	print("\033[?1049l")
	terminal.Restore()
}

func getTermHeight() int {
	_, h, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalln("Could not get terminal size")
	}
	return h
}

func clearScreen() {
	if os.Getenv("NOCLEAR") == "1" {
		return
	}
	print("\033[H\033[2J")
}

func moveCursor(line, col int) {
	fmt.Printf("\033[%d;%dH", line, col)
}
