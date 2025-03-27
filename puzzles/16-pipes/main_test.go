package main

import (
	"testing"
)

func Test_Pipe(t *testing.T) {
	t.Run("make a pipe and rotate", func(t *testing.T) {
		p := NewPipe('╣')

		if p.edges.Up() != 2 {
			t.Errorf("invalid up edges")
		}

		if p.edges.Down() != 2 {
			t.Errorf("invalid down edges")
		}

		if p.edges.Left() != 2 {
			t.Errorf("invalid left edges")
		}

		if p.edges.Right() != 0 {
			t.Errorf("invalid right edges")
		}

		p.Rotate()

		if p.edges.Up() != 2 {
			t.Errorf("invalid up edges")
		}

		if p.edges.Down() != 0 {
			t.Errorf("invalid down edges")
		}

		if p.edges.Left() != 2 {
			t.Errorf("invalid left edges")
		}

		if p.edges.Right() != 2 {
			t.Errorf("invalid right edges")
		}

		if p.char != '╩' {
			t.Errorf("invalid char")
		}
	})

	t.Run("rotatability", func(t *testing.T) {
		s := NewPipe(' ')
		if !s.locked {
			t.Errorf("spaces should be locked")
		}

		x := NewPipe('╬')
		if !x.locked {
			t.Errorf("double cross should be locked")
		}

		sx := NewPipe('┼')
		if !sx.locked {
			t.Errorf("single cross should be locked")
		}

		l := NewPipe('╧')
		if l.locked {
			t.Errorf("╧ should not be locked")
		}
	})

	t.Run("test matches", func(t *testing.T) {
		e := runeMap['│']

		u, r, d, l := runeMap['┐'], runeMap[' '], runeMap['┼'], runeMap['╣']
		if !e.Matches(u, r, d, l) {
			t.Errorf("should match")
		}
	})

	t.Run("test potential edge counts", func(t *testing.T) {
		p := NewPipe('─')

		edges := p.edges.ValidEdgeCounts()
		if len(edges) != 2 {
			t.Errorf("expected 2 valid edge counts, not %d", len(p.edges.ValidEdgeCounts()))
		}
	})

	t.Run("edges test", func(t *testing.T) {
		p := NewPipe('└').edges

		if p.Up() != 1 {
			t.Errorf("expected 1 up got %d", p.Up())
		}
		if p.Down() != 0 {
			t.Errorf("expected 0 down got %d", p.Up())
		}
		if p.Left() != 0 {
			t.Errorf("expected 0 left got %d", p.Up())
		}
		if p.Right() != 1 {
			t.Errorf("expected 1 right got %d", p.Up())
		}
	})
}
