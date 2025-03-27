package main

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/lthummus/i18n-puzzles/input"
)

const (
	MaxSMSBytes   = 160
	MaxTweetChars = 140
)

func getCost(x string) int {
	byteCount := len(x)
	runeCount := utf8.RuneCount([]byte(x))

	canTweet := runeCount <= MaxTweetChars
	canSMS := byteCount <= MaxSMSBytes

	if canTweet && canSMS {
		return 13
	} else if canTweet {
		return 7
	} else if canSMS {
		return 11
	} else {
		return 0
	}
}

func main() {
	input, err := input.GetInputLinesUTF8(context.Background(), 1, input.RealInput)
	if err != nil {
		panic(err)
	}

	totalCost := 0

	for _, curr := range input {
		totalCost += getCost(curr)
	}

	fmt.Printf("%d\n", totalCost)
}
