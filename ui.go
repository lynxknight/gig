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

const (
	U_HEADER = "============"

	T_CURSOR_HIDE               = "\033[?25l"
	T_CURSOR_SHOW               = "\033[?25h"
	T_CLEAR_LINE                = "\033[K"
	T_CLEAR_SCREEN              = "\033[H\033[2J"
	T_NEWLINE                   = "\r\n"
	T_CREATE_ALTERNATIVE_SCREEN = "\033[?1049h"
	T_CLOSE_ALTERNATIVE_SCREEN  = "\033[?1049l"
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

type lineBuf struct {
	lines          []string
	dataLinesCount int
}

func (lb *lineBuf) Append(line string) {
	lb.lines = append(lb.lines, line)
}

func (lb *lineBuf) ExtendText(newLines []string) {
	height := getTermHeight()
	for i := 0; i < min(height, len(newLines)); i++ {
		lb.Append(T_CLEAR_LINE)
		lb.Append(newLines[i])
		lb.dataLinesCount++
		if i < height-1 {
			lb.Append(T_NEWLINE)
		}
	}
}

func (lb *lineBuf) FillWithEmptyLines() {
	height := getTermHeight()
	// it feels like sometimes we are clearing too much :D
	fillerLinesCount := height - lb.dataLinesCount
	for i := 0; i < fillerLinesCount; i++ {
		lb.Append(T_CLEAR_LINE)
		if i < fillerLinesCount-1 {
			lb.Append(T_NEWLINE)
		}
	}
}

func (lb *lineBuf) Draw() {
	fmt.Print(strings.Join(lb.lines, ""))
}

func drawUI(branches []branch, query string, cursorpos int) int {
	moveCursor(0, 0)
	lb := lineBuf{}
	lb.Append(T_CURSOR_HIDE)
	// stringsToDisplay contains querystring and header
	stringsToDisplay := displayBranches(query, branches, cursorpos)
	lb.ExtendText(stringsToDisplay)
	lb.FillWithEmptyLines()
	lb.Append(T_CURSOR_SHOW)
	lb.Draw()
	return len(stringsToDisplay) - 2
}

var highlighter = color.New(color.BgWhite, color.FgBlack).SprintfFunc()
var underliner = color.New(color.Underline).SprintfFunc()

func displayBranches(query string, branches []branch, hindex int) []string {
	// Probably it should not know about cursor position and height
	branchesToPrint := make([]string, 2, len(branches)+2)
	branchesToPrint[0] = query
	branchesToPrint[1] = U_HEADER
	if len(branches) == 0 {
		return branchesToPrint
	}
	maxDistanceOfInterest := branches[0].costcache[query].Distance + 4
	for i := range branches {
		// TODO: we might create / not create "validation" function in runtime
		// Empty query => no validation
		score := branches[i].costcache[query]
		if query != "" && score.Distance > maxDistanceOfInterest {
			break
		}
		i1, i2 := score.I1, score.I2
		str := branches[i].name
		if hindex == i {
			// Escape sequences cannot be "nested", i.e. we cannot do
			// H...U..ESC...ESC, first escape will cancel out both H and U, so
			// we need to re-highlight part that comes after underline
			// TODO: investigate colorlib, maybe we can do it better
			first_part := highlighter("%v%v", str[:i1], underliner("%v", str[i1:i2]))
			second_part := highlighter("%v", str[i2:])
			branchesToPrint = append(branchesToPrint, first_part+second_part)
		} else {
			branchesToPrint = append(
				branchesToPrint, fmt.Sprintf(
					"%v%v%v", str[:i1], underliner("%v", str[i1:i2]), str[i2:],
				),
			)
		}
	}
	return branchesToPrint
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
	if err != nil {
		log.Fatalln("Failed to open tty device")
	}
	err = term.RawMode(t)
	if err != nil {
		log.Fatalln("Failed to enter raw mode")
	}
	fmt.Print(T_CREATE_ALTERNATIVE_SCREEN)
	return t
}

func restoreTerm(terminal *term.Term) {
	if os.Getenv("NOCLEAR") != "1" {
		fmt.Print(T_CLOSE_ALTERNATIVE_SCREEN)
	}
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
	fmt.Print(T_CLEAR_SCREEN)
}

func moveCursor(line, col int) {
	fmt.Printf("\033[%d;%dH", line, col)
}
