package main

import (
	"context"
	"fmt"
	"unicode"

	"github.com/lthummus/i18n-puzzles/input"
)

const (
	MaxLength = 12
	MinLength = 4
)

func isPasswordValid(pwd string) bool {
	length := 0

	digitFound := false
	upperFound := false
	lowerFound := false
	nonASCIIFound := false

	for _, r := range pwd {
		length++

		if r > unicode.MaxASCII {
			nonASCIIFound = true
		}

		if unicode.IsDigit(r) {
			digitFound = true
		} else if unicode.IsUpper(r) {
			upperFound = true
		} else if unicode.IsLower(r) {
			lowerFound = true
		}
	}

	lengthOK := length >= MinLength && length <= MaxLength

	return lengthOK && digitFound && upperFound && lowerFound && nonASCIIFound
}

func main() {
	in, err := input.GetInputLinesUTF8(context.Background(), 3, input.RealInput)
	if err != nil {
		panic(err)
	}

	validCount := 0
	for _, curr := range in {
		if isPasswordValid(curr) {
			validCount++
		}
	}

	fmt.Printf("%d\n", validCount)
}
