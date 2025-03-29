package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Pipe(t *testing.T) {
	t.Run("make a pipe and rotate", func(t *testing.T) {
		p := NewPipe('╣')

		assert.Equal(t, uint8(2), p.edges.Up())
		assert.Equal(t, uint8(2), p.edges.Down())
		assert.Equal(t, uint8(2), p.edges.Left())
		assert.Equal(t, uint8(0), p.edges.Right())

		p.Rotate()

		assert.Equal(t, uint8(2), p.edges.Up())
		assert.Equal(t, uint8(0), p.edges.Down())
		assert.Equal(t, uint8(2), p.edges.Left())
		assert.Equal(t, uint8(2), p.edges.Right())

		assert.Equal(t, '╩', p.char)
	})

	t.Run("rotatability", func(t *testing.T) {
		s := NewPipe(' ')
		assert.True(t, s.locked)

		x := NewPipe('╬')
		assert.True(t, x.locked)

		sx := NewPipe('┼')
		assert.True(t, sx.locked)

		l := NewPipe('╧')
		assert.False(t, l.locked)
	})

	t.Run("test matches", func(t *testing.T) {
		e := runeMap['│']

		u, r, d, l := runeMap['┐'], runeMap[' '], runeMap['┼'], runeMap['╣']
		assert.True(t, e.Matches(u, r, d, l))
	})

	t.Run("test potential edge counts", func(t *testing.T) {
		p := NewPipe('─')

		assert.Len(t, p.edges.ValidEdgeCounts(), 2)
	})

	t.Run("edges test", func(t *testing.T) {
		p := NewPipe('└').edges

		assert.Equal(t, uint8(1), p.Up())
		assert.Equal(t, uint8(0), p.Down())
		assert.Equal(t, uint8(0), p.Left())
		assert.Equal(t, uint8(1), p.Right())
	})
}
