package main

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/lthummus/i18n-puzzles/input"
)

var odysseusNames = []string{"ΟΔΥΣΣΕΥΣ", "ΟΔΥΣΣΕΩΣ", "ΟΔΥΣΣΕΙ", "ΟΔΥΣΣΕΑ", "ΟΔΥΣΣΕΥ"}

var greekUppercase = []rune("ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩ")

func rotateCharacter(x rune) rune {
	idx := slices.Index(greekUppercase, x)
	if idx == -1 {
		return x
	}

	return greekUppercase[(idx+1)%len(greekUppercase)]
}

func rotateString(x string) string {
	var sb strings.Builder
	for _, curr := range x {
		_, err := sb.WriteRune(rotateCharacter(curr))
		if err != nil {
			panic(err)
		}
	}
	return sb.String()
}

func stringContainsOdysseus(x string) bool {
	for _, curr := range odysseusNames {
		if strings.Contains(x, curr) {
			return true
		}
	}
	return false
}

func odysseusFound(x string) int {
	x = strings.ToUpper(x)
	for r := range len(greekUppercase) {
		if stringContainsOdysseus(x) {
			return r
		}
		x = rotateString(x)
	}
	return 0
}

func main() {
	lines, err := input.GetInputLinesUTF8(context.Background(), 11, input.RealInput)
	if err != nil {
		panic(err)
	}

	total := 0
	for _, curr := range lines {
		total += odysseusFound(curr)
	}

	fmt.Printf("%d\n", total)
}
