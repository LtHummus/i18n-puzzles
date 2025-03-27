package main

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"

	"github.com/lthummus/i18n-puzzles/input"
)

var entryRegex = regexp.MustCompile(`(.*), (.*): (\d+)`)

var englishReplacer = strings.NewReplacer(
	"Æ", "AE",
	"Ø", "O",
)

var swedishReplacer = strings.NewReplacer(
	"Æ", "Ä",
	"Ø", "Ö",
)

type Entry struct {
	LastName  string
	FirstName string
	Phone     string

	englishNormalizedKey []string
	swedishNormalizedKey []string
	dutchNormalizedKey   []string
}

func generateDutchKey(lastName, firstName string) []string {
	var fnb strings.Builder
	var lnb strings.Builder

	ln := norm.NFKD.String(lastName)
	fn := norm.NFKD.String(firstName)

	ln = englishReplacer.Replace(ln)
	fn = englishReplacer.Replace(fn)

	upperFound := false

	for _, curr := range ln {
		if unicode.IsLetter(curr) {
			if unicode.IsUpper(curr) {
				upperFound = true
			}

			if upperFound {
				lnb.WriteRune(unicode.ToUpper(curr))
			}
		}
	}
	for _, curr := range fn {
		if unicode.IsLetter(curr) {
			fnb.WriteRune(unicode.ToUpper(curr))
		}
	}

	return []string{lnb.String(), fnb.String()}
}

func generateSwedishKey(lastName, firstName string) []string {
	var fnb strings.Builder
	var lnb strings.Builder

	ln := norm.NFC.String(lastName)
	fn := norm.NFC.String(firstName)

	ln = swedishReplacer.Replace(ln)
	fn = swedishReplacer.Replace(fn)

	for _, curr := range ln {
		if unicode.IsLetter(curr) {
			lnb.WriteRune(unicode.ToUpper(curr))
		}
	}
	for _, curr := range fn {
		if unicode.IsLetter(curr) {
			fnb.WriteRune(unicode.ToUpper(curr))
		}
	}

	return []string{lnb.String(), fnb.String()}
}

func generateEnglishKey(lastName, firstName string) []string {
	var fnb strings.Builder
	var lnb strings.Builder

	ln := norm.NFKD.String(lastName)
	fn := norm.NFKD.String(firstName)

	ln = englishReplacer.Replace(ln)
	fn = englishReplacer.Replace(fn)

	for _, curr := range ln {
		if unicode.IsLetter(curr) {
			lnb.WriteRune(unicode.ToUpper(curr))
		}
	}
	for _, curr := range fn {
		if unicode.IsLetter(curr) {
			fnb.WriteRune(unicode.ToUpper(curr))
		}
	}

	return []string{lnb.String(), fnb.String()}
}

func NewEntry(x string) *Entry {
	matches := entryRegex.FindStringSubmatch(x)

	return &Entry{
		LastName:  matches[1],
		FirstName: matches[2],
		Phone:     matches[3],

		englishNormalizedKey: generateEnglishKey(matches[1], matches[2]),
		swedishNormalizedKey: generateSwedishKey(matches[1], matches[2]),
		dutchNormalizedKey:   generateDutchKey(matches[1], matches[2]),
	}
}

func middleNumber(entries []*Entry) int {
	mid := len(entries) / 2
	p, err := strconv.Atoi(entries[mid].Phone)
	if err != nil {
		panic(err)
	}
	return p
}

func main() {
	lines, err := input.GetInputLinesUTF8(context.Background(), 12, input.RealInput)
	if err != nil {
		panic(err)
	}

	entries := make([]*Entry, len(lines))
	for i := range lines {
		entries[i] = NewEntry(lines[i])
	}

	englishSorted := make([]*Entry, len(entries))
	copy(englishSorted, entries)

	sort.Slice(englishSorted, func(i, j int) bool {
		if englishSorted[i].englishNormalizedKey[0] != englishSorted[j].englishNormalizedKey[0] {
			return englishSorted[i].englishNormalizedKey[0] < englishSorted[j].englishNormalizedKey[0]
		}
		return englishSorted[i].englishNormalizedKey[1] < englishSorted[j].englishNormalizedKey[1]
	})

	swedishSorted := make([]*Entry, len(entries))
	copy(swedishSorted, entries)

	swedish := collate.New(language.Swedish)

	sort.Slice(swedishSorted, func(i, j int) bool {
		if swedish.CompareString(swedishSorted[i].swedishNormalizedKey[0], swedishSorted[j].swedishNormalizedKey[0]) != 0 {
			return swedish.CompareString(swedishSorted[i].swedishNormalizedKey[0], swedishSorted[j].swedishNormalizedKey[0]) < 0
		}
		return swedish.CompareString(swedishSorted[i].swedishNormalizedKey[1], swedishSorted[j].swedishNormalizedKey[1]) < 0
	})

	dutchSorted := make([]*Entry, len(entries))
	copy(dutchSorted, entries)

	sort.Slice(dutchSorted, func(i, j int) bool {
		if dutchSorted[i].dutchNormalizedKey[0] != dutchSorted[j].dutchNormalizedKey[0] {
			return dutchSorted[i].dutchNormalizedKey[0] < dutchSorted[j].dutchNormalizedKey[0]
		}
		return dutchSorted[i].dutchNormalizedKey[1] < dutchSorted[j].dutchNormalizedKey[1]
	})

	ans := middleNumber(englishSorted) * middleNumber(swedishSorted) * middleNumber(dutchSorted)

	fmt.Printf("%d\n", ans)
}
