package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	unidecode "golang.org/x/text/encoding/unicode"

	"github.com/lthummus/i18n-puzzles/input"
)

var (
	UTF8ByteOrderMark              = []byte{0xEF, 0xBB, 0xBF}
	UTF16BigEndianByteOrderMark    = []byte{0xFE, 0xFF}
	UTF16LittleEndianByteOrderMark = []byte{0xFF, 0xFE}
)

// all possible decoders
var decoders = []*encoding.Decoder{
	charmap.ISO8859_1.NewDecoder(),
	unidecode.UTF16(unidecode.BigEndian, unidecode.IgnoreBOM).NewDecoder(),
	unidecode.UTF16(unidecode.LittleEndian, unidecode.IgnoreBOM).NewDecoder(),
}

func validBytes(b []byte) bool {
	if bytes.ContainsRune(b, utf8.RuneError) {
		return false
	}

	if !utf8.Valid(b) {
		return false
	}

	s := string(b)

	for _, curr := range s {
		if !unicode.IsLetter(curr) {
			return false
		}
	}

	return true
}

func decodeHex(in string) []string {
	b, err := hex.DecodeString(in)
	if err != nil {
		panic(err)
	}

	// detect and handle strings with byte order marks. If we have a byte order mark, we know exactly what we're dealing
	// with
	if bytes.HasPrefix(b, UTF8ByteOrderMark) {
		return []string{string(b[3:])}
	}

	if bytes.HasPrefix(b, UTF16BigEndianByteOrderMark) {
		dec := unidecode.UTF16(unidecode.BigEndian, unidecode.ExpectBOM).NewDecoder()
		d, err := dec.Bytes(b)
		if err != nil {
			panic(err)
		}
		return []string{string(d)}
	}

	if bytes.HasPrefix(b, UTF16LittleEndianByteOrderMark) {
		dec := unidecode.UTF16(unidecode.LittleEndian, unidecode.ExpectBOM).NewDecoder()
		d, err := dec.Bytes(b)
		if err != nil {
			panic(err)
		}
		return []string{string(d)}
	}

	// if we DON'T have a byte order mark, then that means we have to just try everything and go on ~vibes~
	var ret []string

	if validBytes(b) {
		ret = append(ret, string(b))
	}

	for _, decoder := range decoders {
		d, err := decoder.Bytes(b)
		if err != nil {
			continue
		}

		if !validBytes(d) {
			continue
		}

		ret = append(ret, string(d))
	}

	return ret
}

func main() {
	in, err := input.GetInputUTF8(context.Background(), 13, input.RealInput)
	if err != nil {
		panic(err)
	}
	
	parts := strings.Split(in, "\n\n")

	crosswordInputs := strings.Split(parts[1], "\n")
	crosswords := make([]*regexp.Regexp, len(crosswordInputs))
	for i := range crosswordInputs {
		crosswords[i] = regexp.MustCompile(fmt.Sprintf("^%s$", strings.TrimSpace(crosswordInputs[i])))
	}

	words := map[string]int{}

	lines := strings.Split(parts[0], "\n")
	for i := range lines {
		potentials := decodeHex(lines[i])
		for _, curr := range potentials {
			words[curr] = i + 1
		}
	}

	total := 0
	for _, curr := range crosswords {
		for word, line := range words {
			if curr.MatchString(word) {
				r := curr.String()
				fmt.Printf("%s (%d) matches %s\n", word, line, r[1:len(r)-1])
				total += line
			}
		}
	}

	fmt.Printf("%d\n", total)

}
