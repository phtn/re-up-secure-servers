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

func StripANSI(input string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegex.ReplaceAllString(input, "")
}
func ExpandTabs(input string, tabWidth int) string {
	return strings.ReplaceAll(input, "\t", strings.Repeat(" ", tabWidth))
}
func CountNewLines(input string) int {
	return strings.Count(input, "\n")
}

func dim(s string) string {
	return Gray(s, 0)
}

var (
	tr = dim("╮")
	tl = dim("╭")
	br = dim("╯")
	bl = dim("╰")
	hl = dim("─")
	vl = dim("│")
	sp = dim(" ")
)

func calcX(input []string) (map[int]interface{}, string) {

	var c int
	var widths []int
	var contents []string

	for _, v := range input {
		contents = append(contents, v)
	}

	for _, u := range input {
		cleaned := StripANSI(u)
		expanded := ExpandTabs(cleaned, 5)
		widths = append(widths, len(expanded))
		c += utf8.RuneCountInString(expanded)

	}

	mint := make(map[int]interface{})
	mint[8] = contents[0]
	mint[4] = contents[1]
	mint[widths[2]] = contents[2] + sp

	// for i, j := range widths {
	// 	mint[j] = contents[i]
	// }

	pad := 1
	w := c + (4 * pad)
	x_border := repeat(hl, w)
	return mint, x_border
}

func OpenRect(inputs []string) string {

	m, x := calcX(inputs)

	// var order []string
	// for i, j := range m {
	// 	if i == 8 {
	// 		order[0] = fmt.Sprintf("%s", j)
	// 	}
	// 	if i == 4 {
	// 		order[1] = fmt.Sprintf("%s", j)
	// 	}
	// 	order[2] = fmt.Sprintf("%s", j)
	// }

	content := vl
	for _, v := range m {
		content += fmt.Sprintf(" %s", v)
	}

	top_b := tl + x + tr
	bot_b := bl + x + br

	var result string
	result += top_b + "\n"
	result += content + vl + "\n"
	result += bot_b
	return result
}

func RectII() {
	slice := []string{"there she goes", "there she goes again."}
	for i := range slice {
		fmt.Println(i)
	}
}

// repeat creates a string by repeating a character n times
func repeat(char string, count int) string {
	return fmt.Sprintf("%s", strings.Repeat(string([]rune(char)), count))
}
