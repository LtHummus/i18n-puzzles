package main

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/lthummus/i18n-puzzles/input"
)

type Holiday struct {
	Year  int
	Month time.Month
	Day   int
}

func (h *Holiday) IsDate(when time.Time) bool {
	return h.Year == when.Year() && h.Month == when.Month() && h.Day == when.Day()
}

func (h *Holiday) String() string {
	return fmt.Sprintf("%04d %s %02d", h.Year, h.Month, h.Day)
}

type TOPlapOffice struct {
	Name     string
	TimeZone *time.Location
	Holidays []Holiday
}

func (t *TOPlapOffice) IsOpen(when time.Time) bool {
	officeTime := when.In(t.TimeZone)

	// is it a holiday?
	for _, curr := range t.Holidays {
		if curr.IsDate(officeTime) {
			return false
		}
	}

	officeWeekday := officeTime.Weekday()
	if officeWeekday == time.Saturday || officeWeekday == time.Sunday {
		return false
	}

	openingTime := time.Date(officeTime.Year(), officeTime.Month(), officeTime.Day(), 8, 29, 0, 0, officeTime.Location())
	closingTime := time.Date(officeTime.Year(), officeTime.Month(), officeTime.Day(), 17, 00, 0, 0, officeTime.Location())

	return officeTime.After(openingTime) && officeTime.Before(closingTime)
}

func (t *TOPlapOffice) String() string {
	holidayStrings := make([]string, len(t.Holidays))
	for i := range t.Holidays {
		holidayStrings[i] = t.Holidays[i].String()
	}

	return fmt.Sprintf("%s (TZ = %s). Holidays = %s", t.Name, t.TimeZone.String(), strings.Join(holidayStrings, ","))
}

func anyOfficeOpen(offices []*TOPlapOffice, when time.Time) bool {
	for _, curr := range offices {
		if curr.IsOpen(when) {
			return true
		}
	}
	return false
}

func isInOfficeHolidays(holidays []Holiday, when time.Time) bool {
	for _, curr := range holidays {
		if curr.IsDate(when) {
			return true
		}
	}

	return false
}

func overtimeNeeded(offices []*TOPlapOffice, x string) int {
	// surely you will not regret iterating minute-by-minute for the whole year
	//
	// "Technically, it's O(1) because 2022 is always 525600 minutes!" -- Me, trying to justify my bad decisions

	overtimeMinutes := 0

	fields := strings.Split(x, "\t")
	officeZone, err := time.LoadLocation(fields[1])
	if err != nil {
		panic(err)
	}
	officeHolidays := decodeHolidays(fields[2])

	currTime := time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2022, time.December, 31, 23, 59, 0, 0, time.UTC)

	minsChecked := 0
	for currTime.Before(endTime) {
		officeTime := currTime.In(officeZone)
		minsChecked++
		if isInOfficeHolidays(officeHolidays, officeTime) {
			currTime = currTime.Add(1 * time.Minute)
			continue
		}

		officeWeekday := officeTime.Weekday()
		if officeWeekday == time.Saturday || officeWeekday == time.Sunday {
			currTime = currTime.Add(1 * time.Minute)
			continue
		}

		if !anyOfficeOpen(offices, currTime) {
			overtimeMinutes++
		}

		currTime = currTime.Add(1 * time.Minute)
	}

	return overtimeMinutes
}

func decodeHolidays(x string) []Holiday {
	unparsedHolidays := strings.Split(x, ";")
	holidays := make([]Holiday, len(unparsedHolidays))

	for i := range unparsedHolidays {
		date, err := time.Parse("2 January 2006", unparsedHolidays[i])
		if err != nil {
			panic(err)
		}
		holidays[i] = Holiday{
			Year:  date.Year(),
			Month: date.Month(),
			Day:   date.Day(),
		}
	}

	return holidays
}

func NewTOPLapOffice(x string) *TOPlapOffice {
	fields := strings.Split(x, "\t")
	name := fields[0]
	timeZone, err := time.LoadLocation(fields[1])
	if err != nil {
		panic(err)
	}

	return &TOPlapOffice{
		Name:     name,
		TimeZone: timeZone,
		Holidays: decodeHolidays(fields[2]),
	}
}

func main() {
	in, err := input.GetInputUTF8(context.Background(), 15, input.RealInput)
	if err != nil {
		panic(err)
	}

	start := time.Now()
	parts := strings.Split(in, "\n\n")

	var toplapOffices []*TOPlapOffice

	for _, curr := range strings.Split(parts[0], "\n") {
		toplapOffices = append(toplapOffices, NewTOPLapOffice(curr))
	}

	customerOfficeLines := strings.Split(parts[1], "\n")

	var wg sync.WaitGroup
	wg.Add(len(customerOfficeLines))

	overtimeOffices := make([]int, len(customerOfficeLines))
	for i := range customerOfficeLines {
		go func() {
			overtimeOffices[i] = overtimeNeeded(toplapOffices, customerOfficeLines[i])
			wg.Done()
		}()
	}

	wg.Wait()

	minReqd := slices.Min(overtimeOffices)
	maxReqd := slices.Max(overtimeOffices)

	fmt.Printf("%d\n", maxReqd-minReqd)
	dur := time.Since(start)
	fmt.Printf("Took %.2f seconds\n", dur.Seconds())
}
