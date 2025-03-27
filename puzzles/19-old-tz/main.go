package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"embed"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/lthummus/i18n-puzzles/input"
)

const (
	TimeFormatString = "2006-01-02 15:04:05"
)

var (
	zoneinfoRegex = regexp.MustCompile(`^zoneinfo-(\d{4}[a-z])\.tar\.gz$`)
)

//go:embed tzdata/*
var tzdata embed.FS

type TzData map[string][]byte

var zonedata = map[string]TzData{}

func loadZonesFromTarGz(x []byte) TzData {
	gzr, err := gzip.NewReader(bytes.NewReader(x))
	if err != nil {
		panic(err)
	}

	ret := TzData{}

	links := map[string]string{}

	tr := tar.NewReader(gzr)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		if strings.Contains(h.Name, ".") {
			continue
		}

		if h.Typeflag == tar.TypeReg {
			var byf bytes.Buffer
			_, err = io.Copy(&byf, tr)
			if err != nil {
				panic(err)
			}
			ret[h.Name] = byf.Bytes()
		} else if h.Typeflag == tar.TypeLink {
			// some entries are symlinked, so save those for processing later
			links[h.Name] = h.Linkname
		}
	}

	for k, v := range links {
		d := make([]byte, len(ret[v]))
		copy(d, ret[v])
		ret[k] = d
	}

	return ret
}

func init() {
	fmt.Printf("Loading available tz versions....\n")
	dir, err := tzdata.ReadDir("tzdata")
	if err != nil {
		panic(err)
	}
	for _, curr := range dir {
		m := zoneinfoRegex.FindStringSubmatch(curr.Name())
		if m != nil {
			fmt.Printf("Loading tzdata %s...", m[1])

			f, err := tzdata.ReadFile(fmt.Sprintf("tzdata/%s", curr.Name()))
			if err != nil {
				panic(err)
			}

			zonedata[m[1]] = loadZonesFromTarGz(f)
			fmt.Printf("DONE!\n")
		}
	}
}

func allContain(x map[string][]int64, v int64) bool {
	for _, c := range x {
		if !slices.Contains(c, v) {
			return false
		}
	}

	return true
}

func main() {
	in, err := input.GetInputLinesUTF8(context.Background(), 19, input.RealInput)
	if err != nil {
		panic(err)
	}

	stationsMoments := map[string][]int64{}
	var firstZone *string

	for _, curr := range in {
		parts := strings.Split(curr, "; ")

		t := parts[0]
		zone := parts[1]

		if firstZone == nil {
			firstZone = &zone
		}

		for k, v := range zonedata {
			l, err := time.LoadLocationFromTZData(zone, v[zone])
			if err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: skipping input line %s because I couldn't load zone %s from version %s\n", curr, zone, k)
				continue
			}

			moment, err := time.ParseInLocation(TimeFormatString, t, l)
			if err != nil {
				panic(err)
			}

			utcTime := moment.UTC().Unix()
			m := stationsMoments[zone]
			m = append(m, utcTime)
			stationsMoments[zone] = m
		}
	}

	if firstZone == nil {
		panic("didn't see any zones?")
	}

	toCheck := stationsMoments[*firstZone]
	for _, curr := range toCheck {
		if allContain(stationsMoments, curr) {
			winner := time.Unix(curr, 0).UTC()
			fmt.Printf("%d -- %s\n", curr, winner.Format("2006-01-02T15:04:05-07:00"))
			break
		}
	}

}
