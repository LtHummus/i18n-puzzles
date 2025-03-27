package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lthummus/i18n-puzzles/input"
)

const (
	TimeFormatString       = "2006-01-02T15:04:05-07:00"
	DetectionCountRequired = 4
)

func main() {
	input, err := input.GetInputLinesUTF8(context.Background(), 2, input.RealInput)
	if err != nil {
		panic(err)
	}

	seenTimes := map[time.Time]int{}

	for i := range input {
		t, err := time.Parse(TimeFormatString, input[i])
		if err != nil {
			panic(err)
		}
		t = t.In(time.UTC)
		seen := seenTimes[t]
		seenTimes[t] = seen + 1
	}

	var detectionTime time.Time
	var found bool

	for k, v := range seenTimes {
		if v >= DetectionCountRequired {
			detectionTime = k
			found = true
		}
	}

	if !found {
		panic("no detection found")
	}

	fmt.Printf("%s\n", detectionTime.Format(TimeFormatString))

}
