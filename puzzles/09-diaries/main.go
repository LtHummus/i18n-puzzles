package main

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/lthummus/i18n-puzzles/input"
)

var (
	Orderings = map[string]string{
		"YMD": "06-01-02",
		"YDM": "06-02-01",
		"MDY": "01-02-06",
		"DMY": "02-01-06",
	}

	MaxDate = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Add(1 * time.Second)
	MinDate = time.Date(1920, 1, 1, 0, 0, 0, 0, time.UTC).Add(-1 * time.Second)

	LineRegex = regexp.MustCompile(`^(\d{2}-\d{2}-\d{2}): (.*)$`)
)

func validOrderings(date string) []string {
	var ret []string
	for ord, form := range Orderings {
		d, err := time.Parse(form, date)
		if err != nil {
			continue
		}

		// go uses a pivot year of 1969 and we want one for 1920
		if d.Year() >= 2020 && (d.Day() != 01 || d.Month() != time.January) {
			d = d.AddDate(-100, 0, 0)
		}

		if d.Before(MinDate) {
			continue
		}

		if d.After(MaxDate) {
			continue
		}

		ret = append(ret, ord)
	}

	return ret
}

func determineOrder(dates []string) string {
	potentialOrderings := validOrderings(dates[0])

	idx := 1

	for len(potentialOrderings) > 1 && idx < len(dates) {
		newOrderings := validOrderings(dates[idx])

		var newPotentials []string

		for _, curr := range potentialOrderings {
			if slices.Contains(newOrderings, curr) {
				newPotentials = append(newPotentials, curr)
			}
		}

		potentialOrderings = newPotentials

		idx++
	}

	if len(potentialOrderings) > 1 {
		panic("multiple potentials found")
	}

	if len(potentialOrderings) == 0 {
		panic("no potential ordering found")
	}

	return potentialOrderings[0]
}

func has911Entry(ord string, dates []string) bool {
	for _, curr := range dates {
		d, err := time.Parse(Orderings[ord], curr)
		if err != nil {
			continue
		}

		// don't need to fix pivot year here since we are looking for 2001 specifically
		if d.Year() == 2001 && d.Day() == 11 && d.Month() == time.September {
			return true
		}
	}

	return false
}

func main() {
	people := map[string][]string{}

	lines, err := input.GetInputLinesUTF8(context.Background(), 9, input.RealInput)
	if err != nil {
		panic(err)
	}

	for _, curr := range lines {
		matches := LineRegex.FindStringSubmatch(curr)
		if matches == nil {
			panic("line does not match")
		}

		date := matches[1]

		names := strings.Split(matches[2], ", ")

		for _, person := range names {
			entries := people[person]
			entries = append(entries, date)
			people[person] = entries
		}
	}

	var wroteAbout911 []string

	for person, entries := range people {
		ord := determineOrder(entries)

		if has911Entry(ord, entries) {
			wroteAbout911 = append(wroteAbout911, person)
		}
	}

	slices.Sort(wroteAbout911)

	fmt.Printf("%s\n", strings.Join(wroteAbout911, " "))
}
