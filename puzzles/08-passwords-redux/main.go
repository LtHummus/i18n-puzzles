package main

import (
	"context"
	"fmt"
	"unicode"

	"golang.org/x/text/unicode/norm"

	"github.com/lthummus/i18n-puzzles/input"
)

const (
	MaxLength = 12
	MinLength = 4
)

var (
	vowels = map[rune]bool{
		'a': true,
		'e': true,
		'i': true,
		'o': true,
		'u': true,
	}
)

func isPasswordValid(pwd string) bool {
	length := 0

	digitFound := false
	vowelFound := false
	consFound := false

	seen := map[rune]bool{}

	var iter norm.Iter
	iter.InitString(norm.NFD, pwd)

	for !iter.Done() {
		length++

		r := iter.Next()

		baseRune := rune(r[0])

		if unicode.IsDigit(baseRune) {
			digitFound = true
		} else if unicode.IsLetter(baseRune) {
			baseRune = unicode.ToLower(baseRune)
			if vowels[baseRune] {
				vowelFound = true
			} else {
				consFound = true
			}
		}

		// have to do this AFTER normalization shenanigans
		if seen[baseRune] {
			return false
		}
		seen[baseRune] = true
	}

	lengthOK := length >= MinLength && length <= MaxLength

	return lengthOK && digitFound && consFound && vowelFound
}

func main() {
	lines, err := input.GetInputLinesUTF8(context.Background(), 8, input.RealInput)
	if err != nil {
		panic(err)
	}
	
	validCount := 0
	for _, curr := range lines {
		if isPasswordValid(curr) {
			validCount++
		}
	}

	fmt.Printf("%d\n", validCount)
}
