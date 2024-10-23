package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func drawBorder(l int, t int) {

	rt := []string{"╭", "╮"}
	rb := []string{"╰", "╯"}

	c := rt
	if t == 1 {
		c = rb
	}

	fmt.Printf(Gray(c[0], 0))
	h_l := "──"
	for range l {
		time.Sleep(5 * time.Millisecond)
		fmt.Printf(Gray(h_l, 0))
	}
	fmt.Println(Gray(c[1], 0))
}

func renderContent(c string, l int) {

	p := l - countVis(c)
	if p < 0 {
		p = p * -1
	}

	v := Gray("│", 0)
	ws := strings.Repeat(" ", l-countVis(c))

	// fmt.Printf(" %v ", (l-countVis(c))/2)
	fmt.Println(v + c + ws + v + " " + strconv.Itoa(countVis(c)))
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func countVis(input string) int {
	// Regular expression to match ANSI escape sequences
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	// Remove all ANSI escape sequences from the input string
	cleanString := ansiRegex.ReplaceAllString(input, "")
	// Count the number of characters in the cleaned string
	return utf8.RuneCountInString(cleanString)
}

func MkOne() {

	clearScreen()

	f0 := Light("   ", 0)
	s0 := Light(" ", 0)
	// f1 := Light("    ⟢      ", 0) + DSky("╭", 0)
	f2 := DSky("╭", 0) + Light(" ╮  ", 0)
	f := f0 + f2
	// s1 := Light(" ⟢       ", 0) + DSky("╭", 0) + Slate("◜", 0)
	s2 := DSky("╭", 0) + Light("◜", 0) + DSky("╰", 0) + Light(" ╯", 0) + Gray("◝𐑪", 0)
	s := s0 + s2
	// q := Light("   ⟢        ", 0) + Sky("◌ ", 0)
	t := f0 + " " + Sky("◌   ", 0)

	drawBorder(5, 0)
	for c := range map[string]interface{}{f: f, s: s, t: t} {
		renderContent(c, 10)
	}
	drawBorder(5, 1)
}
