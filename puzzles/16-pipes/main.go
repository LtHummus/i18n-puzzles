package main

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"golang.org/x/text/encoding/charmap"

	"github.com/lthummus/i18n-puzzles/input"
)

const (
	Up = iota * 2
	Right
	Down
	Left
)

const (
	RealFrameLeftEdge          = " │ ║   "
	RealFrameRightEdge         = "   ║ │░"
	RealFrameTopBottomEdgeSize = 5
)

var runeMap = map[rune]Edges{
	' ': edgesFromAdj(0, 0, 0, 0),
	'│': edgesFromAdj(1, 0, 1, 0),
	'┤': edgesFromAdj(1, 0, 1, 1),
	'╡': edgesFromAdj(1, 0, 1, 2),
	'╢': edgesFromAdj(2, 0, 2, 1),
	'╖': edgesFromAdj(0, 0, 2, 1),
	'╕': edgesFromAdj(0, 0, 1, 2),
	'╣': edgesFromAdj(2, 0, 2, 2),
	'║': edgesFromAdj(2, 0, 2, 0),
	'╗': edgesFromAdj(0, 0, 2, 2),
	'╝': edgesFromAdj(2, 0, 0, 2),
	'╜': edgesFromAdj(2, 0, 0, 1),
	'╛': edgesFromAdj(1, 0, 0, 2),
	'┐': edgesFromAdj(0, 0, 1, 1),
	'└': edgesFromAdj(1, 1, 0, 0),
	'┴': edgesFromAdj(1, 1, 0, 1),
	'┬': edgesFromAdj(0, 1, 1, 1),
	'├': edgesFromAdj(1, 1, 1, 0),
	'─': edgesFromAdj(0, 1, 0, 1),
	'┼': edgesFromAdj(1, 1, 1, 1),
	'╞': edgesFromAdj(1, 2, 1, 0),
	'╟': edgesFromAdj(2, 1, 2, 0),
	'╚': edgesFromAdj(2, 2, 0, 0),
	'╔': edgesFromAdj(0, 2, 2, 0),
	'╩': edgesFromAdj(2, 2, 0, 2),
	'╦': edgesFromAdj(0, 2, 2, 2),
	'╠': edgesFromAdj(2, 2, 2, 0),
	'═': edgesFromAdj(0, 2, 0, 2),
	'╬': edgesFromAdj(2, 2, 2, 2),
	'╧': edgesFromAdj(1, 2, 0, 2),
	'╨': edgesFromAdj(2, 1, 0, 1),
	'╤': edgesFromAdj(0, 2, 1, 2),
	'╥': edgesFromAdj(0, 1, 2, 1),
	'╙': edgesFromAdj(2, 1, 0, 0),
	'╘': edgesFromAdj(1, 2, 0, 0),
	'╒': edgesFromAdj(0, 2, 1, 0),
	'╓': edgesFromAdj(0, 1, 2, 0),
	'╫': edgesFromAdj(2, 1, 2, 1),
	'╪': edgesFromAdj(1, 2, 1, 2),
	'┘': edgesFromAdj(1, 0, 0, 1),
	'┌': edgesFromAdj(0, 1, 1, 0),
}
var edgesMap = map[Edges]rune{}

func init() {
	for k, v := range runeMap {
		edgesMap[v] = k
	}
}

const directionMask = 0b11

type Edges uint8

func edgesFromAdj(u, r, d, l uint8) Edges {
	return Edges((u << Up) + (r << Right) + (d << Down) + (l << Left))
}

func (e Edges) Up() uint8 {
	return uint8((e >> Up) & directionMask)
}

func (e Edges) Right() uint8 {
	return uint8((e >> Right) & directionMask)
}

func (e Edges) Down() uint8 {
	return uint8((e >> Down) & directionMask)
}

func (e Edges) Left() uint8 {
	return uint8((e >> Left) & directionMask)
}

func (e Edges) Rotate() Edges {
	return Edges(uint8(e)<<2 + e.Left())
}

func (e Edges) Rotatable() bool {
	return e.Rotate() != e
}

func (e Edges) Matches(u, r, d, l Edges) bool {
	return e.Up() == u.Down() && e.Left() == l.Right() && e.Down() == d.Up() && e.Right() == r.Left()
}

func (e Edges) ValidEdgeCounts() []uint8 {
	seen := map[uint8]bool{}
	seen[e.Up()] = true
	seen[e.Right()] = true
	seen[e.Down()] = true
	seen[e.Left()] = true

	var ret []uint8
	for k := range seen {
		ret = append(ret, k)
	}
	return ret
}

func (e Edges) AllRotations() []Edges {
	seen := map[Edges]bool{}

	r1 := e.Rotate()
	r2 := r1.Rotate()
	r3 := r2.Rotate()

	seen[e] = true
	seen[r1] = true
	seen[r2] = true
	seen[r3] = true

	var ret []Edges
	for k := range seen {
		ret = append(ret, k)
	}

	return ret
}

func (e Edges) ValidRotations(u, r, d, l *Pipe) []Edges {
	potentials := e.AllRotations()

	var potentialUps []uint8
	var potentialDowns []uint8
	var potentialLefts []uint8
	var potentialRights []uint8

	if u.locked {
		potentialUps = []uint8{u.edges.Down()}
	} else {
		potentialUps = u.edges.ValidEdgeCounts()
	}

	if d.locked {
		potentialDowns = []uint8{d.edges.Up()}
	} else {
		potentialDowns = d.edges.ValidEdgeCounts()
	}

	if l.locked {
		potentialLefts = []uint8{l.edges.Right()}
	} else {
		potentialLefts = l.edges.ValidEdgeCounts()
	}

	if r.locked {
		potentialRights = []uint8{r.edges.Left()}
	} else {
		potentialRights = r.edges.ValidEdgeCounts()
	}

	var ret []Edges

	for _, curr := range potentials {
		if slices.Contains(potentialUps, curr.Up()) &&
			slices.Contains(potentialDowns, curr.Down()) &&
			slices.Contains(potentialLefts, curr.Left()) &&
			slices.Contains(potentialRights, curr.Right()) {

			ret = append(ret, curr)
		}
	}

	return ret
}

type Pipe struct {
	char   rune
	edges  Edges
	locked bool
}

func NewPipe(x rune) Pipe {
	e := runeMap[x]
	c := edgesMap[e] // get rid of things we don't care about

	return Pipe{
		char:   c,
		edges:  e,
		locked: !e.Rotatable(),
	}
}

func (p *Pipe) Rotate() {
	if p.locked {
		panic("can not rotate locked pipe")
	}

	p.edges = p.edges.Rotate()

	p.char = edgesMap[p.edges]
}

func (p *Pipe) String() string {
	return fmt.Sprintf("%c", p.char)
}

type Maze struct {
	startX int
	startY int

	endX int
	endY int

	rotations int

	pipes [][]Pipe
}

func NewMaze(x []byte) *Maze {
	mazeBytes, err := charmap.CodePage437.NewDecoder().Bytes(x)
	if err != nil {
		panic(err)
	}

	mazeString := string(mazeBytes)
	lines := strings.Split(mazeString, "\r\n")

	lines = lines[:len(lines)-1]

	if len(lines) != 8 {
		// real input, so modify some things

		// first remove top and bottom frames
		lines = lines[RealFrameTopBottomEdgeSize : len(lines)-RealFrameTopBottomEdgeSize]

		for i, curr := range lines {
			curr, _ = strings.CutPrefix(curr, RealFrameLeftEdge)
			curr, _ = strings.CutSuffix(curr, RealFrameRightEdge)
			lines[i] = curr
		}
	}

	var pipes [][]Pipe
	for _, currLine := range lines {
		var linePipe []Pipe
		for _, c := range currLine {
			linePipe = append(linePipe, NewPipe(c))
		}
		pipes = append(pipes, linePipe)
	}

	return &Maze{
		pipes: pipes,

		startX: 0,
		startY: 0,
		endX:   len(pipes[0]) - 1,
		endY:   len(pipes) - 1,
	}
}

func (m *Maze) getPipe(x, y int) *Pipe {
	if x == 0 && y == -1 {
		// special case, this is a locked '|' character
		p := NewPipe('│')
		p.locked = true
		return &p
	} else if x == m.endX && y == m.endY+1 {
		p := NewPipe('│')
		p.locked = true
		return &p
	}

	if x >= 0 && x < len(m.pipes[0]) && y >= 0 && y < len(m.pipes) {
		return &m.pipes[y][x]
	} else {
		p := NewPipe(' ')
		// set to locked inside NewPipe so don't need to set it here
		return &p
	}
}

func (m *Maze) isSolved() bool {
	for y := range m.pipes {
		for _, curr := range m.pipes[y] {
			if !curr.locked {
				return false
			}
		}
	}

	return true
}

func (m *Maze) recomputeLocked() {
	for y := range m.pipes {
		for x := range m.pipes[y] {
			curr := &m.pipes[y][x]

			if curr.locked {
				continue
			}

			up := m.getPipe(x, y-1)
			right := m.getPipe(x+1, y)
			down := m.getPipe(x, y+1)
			left := m.getPipe(x-1, y)

			potentialRotations := curr.edges.ValidRotations(up, right, down, left)
			if len(potentialRotations) == 0 {
				panic(fmt.Sprintf("no valid rotations at (%d, %d)", x, y))
			} else if len(potentialRotations) == 1 {
				// lock this cell, there's only one way it can go
				ourRotations := 0
				for curr.edges != potentialRotations[0] {
					curr.Rotate()
					ourRotations++
				}
				curr.char = edgesMap[potentialRotations[0]]
				curr.edges = potentialRotations[0]
				curr.locked = true

				m.rotations += ourRotations
				fmt.Printf("Locked %c (%d, %d) after %d rotations (%d total so far)\n", curr.char, x, y, ourRotations, m.rotations)
			}
		}
	}
}

func (m *Maze) String() string {
	var lines []string
	for _, currLine := range m.pipes {
		var sb strings.Builder
		for _, curr := range currLine {
			sb.WriteRune(curr.char)
		}
		lines = append(lines, sb.String())
	}

	return fmt.Sprintf("%s", strings.Join(lines, "\n"))
}

func main() {
	in, err := input.GetInputBytes(context.Background(), 16, input.RealInput)
	if err != nil {
		panic(err)
	}

	m := NewMaze(in)

	fmt.Printf("%s\n", m)

	cycles := 0
	for !m.isSolved() {
		cycles++
		m.recomputeLocked()
	}

	fmt.Printf("%s\n", m)

	fmt.Printf("Solved after %d cycles\n", cycles)
	fmt.Printf("%d\n", m.rotations)
}
