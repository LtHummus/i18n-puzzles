package main

import (
	"context"
	"fmt"
	"math"
	"slices"
	"strings"
	"unicode"

	"github.com/lthummus/i18n-puzzles/input"
)

const (
	RLIMarker = '\u2067'
	LRIMarker = '\u2066'
	PDIMarker = '\u2069'
)

func stripMarkers(x string) string {
	var sb strings.Builder
	for _, curr := range x {
		if curr == RLIMarker || curr == LRIMarker || curr == PDIMarker {
			continue
		}
		sb.WriteRune(curr)
	}

	return sb.String()
}

func computeWithoutUnicode(in string) int {
	stripped := stripMarkers(in)
	fmt.Printf("%s\n", stripped)

	return 0
}

func findHighestRun(embeddingLevels []int) (int, int, int) {
	highestLevel := -1
	highestIdx := -1
	highestLength := -1
	inHighest := false

	for i, curr := range embeddingLevels {
		if curr < highestLevel {
			inHighest = false
			continue
		}

		if inHighest && curr == highestLevel {
			highestLength++
			continue
		}

		if curr > highestLevel {
			highestLevel = curr
			highestLength = 1
			highestIdx = i
			inHighest = true
		}
	}

	return highestLevel, highestIdx, highestLength
}

func fixReversedString(in string) string {
	var embeddingLevels []int
	var chars []rune

	currLevel := 0
	for _, c := range in {
		increasedForDigit := false
		if currLevel%2 == 1 && unicode.IsDigit(c) {
			currLevel++
			increasedForDigit = true
		}
		chars = append(chars, c)
		embeddingLevels = append(embeddingLevels, currLevel)
		if c == RLIMarker && currLevel%2 == 0 {
			currLevel++
		} else if c == LRIMarker && currLevel%2 == 1 {
			currLevel++
		} else if c == PDIMarker {
			// overwrite the last one since we wrote alread
			currLevel--
			embeddingLevels[len(embeddingLevels)-1] = currLevel
		}
		if increasedForDigit {
			currLevel--
		}
	}

	finalRunes := []rune(in)

	highestLevel, highestIdx, highestLength := findHighestRun(embeddingLevels)
	for highestLevel != 0 {
		if highestLength == 1 {
			// trivial case
			embeddingLevels[highestIdx]--
			highestLevel, highestIdx, highestLength = findHighestRun(embeddingLevels)
			continue
		}
		runesToReverse := finalRunes[highestIdx : highestIdx+highestLength]
		slices.Reverse(runesToReverse)

		for i, curr := range runesToReverse {
			if curr == '(' {
				runesToReverse[i] = ')'
			} else if curr == ')' {
				runesToReverse[i] = '('
			}
		}

		var newRuneList []rune
		for _, c := range finalRunes[:highestIdx] {
			newRuneList = append(newRuneList, c)
		}
		for _, c := range runesToReverse {
			newRuneList = append(newRuneList, c)
		}
		for _, c := range finalRunes[highestIdx+highestLength:] {
			newRuneList = append(newRuneList, c)
		}

		finalRunes = newRuneList

		for i := highestIdx; i < highestIdx+highestLength; i++ {
			embeddingLevels[i] = embeddingLevels[i] - 1
		}

		highestLevel, highestIdx, highestLength = findHighestRun(embeddingLevels)
	}

	return stripMarkers(string(finalRunes))
}

func main() {
	in, err := input.GetInputLinesUTF8(context.Background(), 18, input.RealInput)
	if err != nil {
		panic(err)
	}

	total := float64(0)

	for _, curr := range in {
		a := stripMarkers(curr)
		b := fixReversedString(curr)

		an := evalExpression(a)
		bn := evalExpression(b)

		total += math.Abs(an - bn)
	}

	fmt.Printf("%.f\n", total)
}
