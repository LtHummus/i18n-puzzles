package main

import (
	"context"
	"fmt"

	"github.com/lthummus/i18n-puzzles/input"
)

const (
	dx = 2
	dy = 1

	poop = '💩'
)

func main() {
	currX := 0
	currY := 0

	lines, err := input.GetInputLinesUTF8(context.Background(), 5, input.RealInput)
	if err != nil {
		panic(err)
	}

	runeLines := make([][]rune, len(lines))

	for i := range lines {
		runeLines[i] = []rune(lines[i])
	}

	height := len(runeLines)
	width := len(runeLines[0])

	fmt.Printf("%dx%d\n", width, height)

	poops := 0

	for currY < height-1 {
		curr := runeLines[currY%height][currX%width]
		if curr == poop {
			poops += 1
		}

		currX += dx
		currY += dy
	}

	fmt.Printf("%d\n", poops)

}
