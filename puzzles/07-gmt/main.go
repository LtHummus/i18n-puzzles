package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lthummus/i18n-puzzles/input"
)

var (
	halifaxZone  *time.Location
	santiagoZone *time.Location
)

func init() {
	var err error
	halifaxZone, err = time.LoadLocation("America/Halifax")
	if err != nil {
		panic(err)
	}

	santiagoZone, err = time.LoadLocation("America/Santiago")
	if err != nil {
		panic(err)
	}
}

func fix(entry string) time.Time {
	parts := strings.Split(entry, "\t")

	parsed, err := time.Parse("2006-01-02T15:04:05.000-07:00", parts[0])
	if err != nil {
		panic(err)
	}

	_, offset := parsed.Zone()
	utc := parsed.UTC()

	halifaxTime := utc.In(halifaxZone)

	_, halifaxOffset := halifaxTime.Zone()

	var fixedTime time.Time
	if halifaxOffset == offset {
		fixedTime = halifaxTime
	} else {
		fixedTime = utc.In(santiagoZone)
	}

	toAdd, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}

	toSub, err := strconv.Atoi(parts[2])
	if err != nil {
		panic(err)
	}

	fixedTime = fixedTime.Add(time.Duration(toAdd) * time.Minute)
	fixedTime = fixedTime.Add(time.Duration(-toSub) * time.Minute)

	return fixedTime

}

func main() {
	in, err := input.GetInputLinesUTF8(context.Background(), 7, input.RealInput)
	if err != nil {
		panic(err)
	}

	var total int

	for i := range in {
		lineNum := i + 1
		fixedTime := fix(in[i])
		total += fixedTime.Hour() * lineNum
	}

	fmt.Printf("%d\n", total)
}
