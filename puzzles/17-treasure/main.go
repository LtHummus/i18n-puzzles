package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/lthummus/i18n-puzzles/input"
)

const (
	theX     = '╳'
	nullByte = byte(0)
)

var (
	LeftEdges = []string{"╔", "|", "║", "╚"}
)

type ByteKind int

const (
	SingleByteCodePoint ByteKind = iota
	TwoByteHeader
	ThreeByteHeader
	FourByteHeader
	ContinuationByte
)

func detectByteKind(x byte) ByteKind {
	if x&0x80 == 0 {
		return SingleByteCodePoint
	}

	if x&0xE0 == 0xC0 {
		return TwoByteHeader
	}

	if x&0xF0 == 0xE0 {
		return ThreeByteHeader
	}

	if x&0xF8 == 0xF0 {
		return FourByteHeader
	}

	if x&0xC0 == 0x80 {
		return ContinuationByte
	}

	panic("unknown byte kind")
}

type Chunk struct {
	hexInput []string
	input    []string
	lines    [][]byte

	bytesMissingAtEnd       []int
	bytesMissingAtBeginning []int

	used bool
}

func (c *Chunk) meshesRightWithEdge(edge []int) bool {
	toCheck := min(len(edge), len(c.bytesMissingAtEnd))
	for i := 0; i < toCheck; i++ {
		if edge[i] != c.bytesMissingAtEnd[i] {
			return false
		}
	}

	return true
}

func (c *Chunk) meshesLeftWithEdge(edge []int) bool {
	//return slices.Equal(c.bytesMissingAtBeginning[:len(edge)], edge)
	toCheck := min(len(edge), len(c.bytesMissingAtBeginning))
	for i := 0; i < toCheck; i++ {
		if edge[i] != c.bytesMissingAtBeginning[i] {
			return false
		}
	}

	return true
}

func (c *Chunk) isTopEdge() bool {
	return strings.Contains(c.input[0], "-═") || strings.Contains(c.input[0], "╔") || strings.Contains(c.input[0], "╗")
}

func (c *Chunk) isLeftEdge() bool {
	for _, curr := range c.input {
		for _, r := range LeftEdges {
			if strings.HasPrefix(curr, r) {
				return true
			}
		}
	}

	return false
}

func detectDanglingContinuationBytes(x []byte) int {
	continuationBytesFound := 0
	for i := range x {
		if detectByteKind(x[i]) == ContinuationByte {
			continuationBytesFound++
		} else {
			break
		}
	}

	return continuationBytesFound
}

func detectEndBytesMissing(x []byte) int {
	continuationBytesFound := 0

	var i int
	for i = len(x) - 1; i >= 0; i-- {
		if detectByteKind(x[i]) != ContinuationByte {
			break
		}
		continuationBytesFound++
	}

	if i < 0 {
		panic("oops all continuation bytes?")
	}

	bk := detectByteKind(x[i])
	switch bk {
	case SingleByteCodePoint:
		if continuationBytesFound > 0 {
			panic("continuation bytes found after single byte header")
		}
		return 0
	case TwoByteHeader:
		return 1 - continuationBytesFound
	case ThreeByteHeader:
		return 2 - continuationBytesFound
	case FourByteHeader:
		return 3 - continuationBytesFound
	case ContinuationByte:
		panic("continuation byte detected where it shouldn't be?")
	default:
		panic("unknown byte kind")
	}

}

func NewChunk(x string) Chunk {
	lines := strings.Split(x, "\n")

	var b [][]byte
	var endBytesMissing []int
	var beginBytesMissing []int

	var rawLines []string

	for _, curr := range lines {
		lineBytes, err := hex.DecodeString(curr)
		if err != nil {
			panic(err)
		}

		if len(lineBytes) == 0 {
			continue
		}

		endEdge := detectEndBytesMissing(lineBytes)
		beginEdge := detectDanglingContinuationBytes(lineBytes)

		b = append(b, lineBytes)
		endBytesMissing = append(endBytesMissing, endEdge)
		beginBytesMissing = append(beginBytesMissing, beginEdge)
		rawLines = append(rawLines, string(lineBytes))
	}

	return Chunk{
		hexInput:                lines,
		input:                   rawLines,
		lines:                   b,
		bytesMissingAtBeginning: beginBytesMissing,
		bytesMissingAtEnd:       endBytesMissing,
	}
}

func puzzleComplete(chunks []*Chunk) bool {
	unusedChunks := 0
	for _, curr := range chunks {
		if !curr.used {
			unusedChunks++
		}
	}
	return unusedChunks == 0
}

func printMap(m [][]byte) {
	s := make([]string, len(m))
	for i := range m {
		s[i] = string(m[i])
	}
	fmt.Printf("%s\n", strings.Join(s, "\n"))
}

func findPuzzleEdge(puzzleLine []byte) (int, bool) {
	for i := 0; i < len(puzzleLine)-1; i++ {
		currByteNull := puzzleLine[i] == nullByte
		nextByteNull := puzzleLine[i+1] == nullByte

		// this is essentially an XOR
		if currByteNull != nextByteNull {
			goingRight := !currByteNull
			trueIndex := i

			// the caller of this expects us to return the position of the first null byte, which is the next one in
			// the case of going left
			if goingRight {
				trueIndex += 1
			}
			return trueIndex, goingRight
		}
	}

	panic("no edge found")
}

func main() {
	in, err := input.GetInputUTF8(context.Background(), 17, input.RealInput)
	if err != nil {
		panic(err)
	}

	start := time.Now()
	chunkInputs := strings.Split(in, "\n\n")

	chunks := make([]*Chunk, len(chunkInputs))
	for i := range chunkInputs {
		c := NewChunk(chunkInputs[i])
		chunks[i] = &c
	}

	fmt.Printf("Found %d chunks\n", len(chunks))

	var upperLeft *Chunk

	for _, curr := range chunks {
		for _, l := range curr.lines {
			if bytes.ContainsRune(l, '╔') {
				upperLeft = curr
			}
		}
	}

	if upperLeft == nil {
		panic("no upper left found")
	}

	height := 0
	width := 0
	for _, curr := range chunks {
		if curr.isLeftEdge() {
			height += len(curr.lines)
		}
		if curr.isTopEdge() {
			width += len(curr.lines[0]) // need number of BYTES not runes
		}
	}

	fmt.Printf("Found height %d\n", height)
	fmt.Printf("Found width %d\n", width)

	var puzzle [][]byte

	for i := 0; i < height; i++ {
		line := make([]byte, width)
		puzzle = append(puzzle, line)
	}

	// seed with upper left corner
	for y, l := range upperLeft.lines {
		for x, c := range l {
			puzzle[y][x] = c
		}
	}
	upperLeft.used = true

	for !puzzleComplete(chunks) {
		// find the top-left most edge that has not been completed
		y := slices.IndexFunc(puzzle, func(i []byte) bool {
			return slices.Contains(i, nullByte)
		})

		// now find the x offset to go
		x, goingRight := findPuzzleEdge(puzzle[y])

		var edges []int
		currY := y
		// now find the height we want
		for {
			if currY > len(puzzle)-1 {
				break
			}
			if goingRight && puzzle[currY][x-1] != nullByte {
				edges = append(edges, detectEndBytesMissing(puzzle[currY][:x]))
				currY++
			} else if !goingRight && puzzle[currY][x+1] != nullByte {
				edges = append(edges, detectDanglingContinuationBytes(puzzle[currY][x+1:]))
				currY++
			} else {
				break
			}
		}

		// now that we have an edge to solve, go find one
		foundChunkIdx := slices.IndexFunc(chunks, func(chunk *Chunk) bool {
			if chunk.used {
				return false
			}
			if goingRight {
				return chunk.meshesLeftWithEdge(edges)
			} else {
				return chunk.meshesRightWithEdge(edges)
			}
		})

		foundChunk := chunks[foundChunkIdx]

		if goingRight {
			for chunkLineNum, chunkLine := range foundChunk.lines {
				for chunkByteIdx, chunkByte := range chunkLine {
					puzzle[y+chunkLineNum][x+chunkByteIdx] = chunkByte
				}
			}
		} else {
			x = x - len(foundChunk.lines[0]) + 1
			for chunkLineNum, chunkLine := range foundChunk.lines {
				for chunkByteIdx, chunkByte := range chunkLine {
					puzzle[y+chunkLineNum][x+chunkByteIdx] = chunkByte
				}
			}
		}

		foundChunk.used = true
	}

	mapStrings := make([]string, len(puzzle))
	for i := range puzzle {
		mapStrings[i] = string(puzzle[i])
	}

	treasureX := slices.IndexFunc(mapStrings, func(s string) bool {
		return strings.ContainsRune(s, theX)
	})

	var treasureY int
	for _, c := range mapStrings[treasureX] {
		if c == theX {
			break
		}
		treasureY++
	}

	dur := time.Since(start)

	printMap(puzzle)

	fmt.Printf("%d\n", treasureY*treasureX)
	fmt.Printf("Took %02dμs\n", dur.Microseconds())
}
