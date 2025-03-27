package main

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/lthummus/i18n-puzzles/input"
)

var (
	entryRegex = regexp.MustCompile(`Departure: ([A-Za-z/_-]+)\s+(.+)\nArrival: {3}([A-Za-z/_-]+)\s+(.+)`)

	timeFormatString = "Jan 02, 2006, 15:04"
)

func flightTime(entryMatches []string) int {
	departZone, err := time.LoadLocation(entryMatches[1])
	if err != nil {
		panic(err)
	}
	arriveZone, err := time.LoadLocation(entryMatches[3])
	if err != nil {
		panic(err)
	}

	departTime, err := time.ParseInLocation(timeFormatString, entryMatches[2], departZone)
	if err != nil {
		panic(err)
	}

	arriveTime, err := time.ParseInLocation(timeFormatString, entryMatches[4], arriveZone)
	if err != nil {
		panic(err)
	}

	return int(arriveTime.Sub(departTime).Minutes())
}

func main() {
	in, err := input.GetInputUTF8(context.Background(), 4, input.RealInput)
	if err != nil {
		panic(err)
	}

	m := entryRegex.FindAllStringSubmatch(in, -1)

	var travelTime int
	for _, curr := range m {
		travelTime += flightTime(curr)
	}

	fmt.Printf("%d\n", travelTime)
}
