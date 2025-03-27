package main

import (
	"context"
	_ "embed"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/charmap"

	"github.com/lthummus/i18n-puzzles/input"
)

func demangle(x string) string {
	ld := charmap.ISO8859_1.NewEncoder()
	x2, err := ld.String(x)
	if err != nil {
		panic(err)
	}

	return x2
}

func findWordSolution(dict []string, template string) int {
	// this is a good idea? lol
	rx := regexp.MustCompile(fmt.Sprintf("^%s$", strings.TrimSpace(template)))

	for i := range dict {
		if rx.MatchString(dict[i]) {
			return i + 1
		}
	}

	panic("no solution")
}

func main() {
	in, err := input.GetInputUTF8(context.Background(), 6, input.RealInput)
	if err != nil {
		panic(err)
	}

	parts := strings.Split(in, "\n\n")
	lines := strings.Split(parts[0], "\n")

	fixedWords := make([]string, len(lines))

	for i := range lines {
		lineNum := i + 1

		curr := lines[i]

		// decode every 3rd and every 5th line. Every 15th line should be decoded twice
		if lineNum%3 == 0 {
			curr = demangle(curr)
		}
		if lineNum%5 == 0 {
			curr = demangle(curr)
		}

		fixedWords[i] = curr
	}

	slots := strings.Split(parts[1], "\n")

	total := 0
	for _, curr := range slots {
		total += findWordSolution(fixedWords, curr)
	}

	fmt.Printf("%d\n", total)
}
