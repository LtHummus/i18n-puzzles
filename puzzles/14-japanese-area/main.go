package main

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/lthummus/i18n-puzzles/input"
)

const (
	delimiter  = " × "
	shakuValue = float64(10) / float64(33)
)

var largeUnits = map[string]int64{
	"尺": 1,
	"間": 6,
	"丈": 10,
	"町": 360,
	"里": 12960,
}

var smallUnits = map[string]float64{
	"毛": 10000,
	"厘": 1000,
	"分": 100,
	"寸": 10,
}

var jaNums = map[rune]int64{
	'一': 1,
	'二': 2,
	'三': 3,
	'四': 4,
	'五': 5,
	'六': 6,
	'七': 7,
	'八': 8,
	'九': 9,
}

var jaTens = map[rune]int64{
	'十': 10,
	'百': 100,
	'千': 1000,
}

var jaMyriads = map[rune]int64{
	'万': 10000,
	'億': 100000000,
}

func convertToMeters(num int64, unit string) float64 {
	if lu, ok := largeUnits[unit]; ok {
		numShaku := num * lu
		return float64(numShaku) * shakuValue
	}

	if su, ok := smallUnits[unit]; ok {
		numShaku := float64(num) / su
		return numShaku * shakuValue
	}

	panic("unknown unit")
}

func parseJapaneseNumber(x string) (int64, error) {
	var total int64

	var myriadRunning int64
	var running int64

	for _, curr := range x {
		if my := jaMyriads[curr]; my != 0 {
			if running != 0 {
				myriadRunning += running
			}
			total += my * myriadRunning
			myriadRunning = 0
			running = 0
		} else if p10 := jaTens[curr]; p10 != 0 {
			if running == 0 {
				running = 1
			}
			myriadRunning += running * p10
			running = 0
		} else if digit := jaNums[curr]; digit != 0 {
			running = digit
		} else {
			return 0, fmt.Errorf("parseJapaneseNumber: %c: invalid japanese digit", curr)
		}
	}

	total += myriadRunning
	total += running

	return total, nil
}

func parseArea(x string) int64 {
	parts := strings.Split(x, delimiter)
	if len(parts) != 2 {
		panic(x)
	}

	aRunes := []rune(parts[0])
	bRunes := []rune(parts[1])

	unitA := string(aRunes[len(aRunes)-1])
	unitB := string(bRunes[len(bRunes)-1])

	numA, err := parseJapaneseNumber(string(aRunes[:len(aRunes)-1]))
	if err != nil {
		panic(err)
	}
	numB, err := parseJapaneseNumber(string(bRunes[:len(bRunes)-1]))
	if err != nil {
		panic(err)
	}

	a := convertToMeters(numA, unitA)
	b := convertToMeters(numB, unitB)

	// problem spec says that each area will be an int
	return int64(math.Round(a * b))
}

func main() {
	lines, err := input.GetInputLinesUTF8(context.Background(), 14, input.RealInput)
	if err != nil {
		panic(err)
	}
	var total int64

	for _, curr := range lines {
		total += parseArea(curr)
	}

	fmt.Printf("%d\n", total)
}
