package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ByteKind(t *testing.T) {
	assert.Equal(t, SingleByteCodePoint, detectByteKind(0x32))

	assert.Equal(t, TwoByteHeader, detectByteKind(0b11011000))
	assert.Equal(t, TwoByteHeader, detectByteKind(0b11011111))
	assert.Equal(t, TwoByteHeader, detectByteKind(0b11000000))

	assert.Equal(t, ThreeByteHeader, detectByteKind(0b11101101))
	assert.Equal(t, ThreeByteHeader, detectByteKind(0b11100000))
	assert.Equal(t, ThreeByteHeader, detectByteKind(0b11101111))

	assert.Equal(t, FourByteHeader, detectByteKind(0b11110101))
	assert.Equal(t, FourByteHeader, detectByteKind(0b11110000))
	assert.Equal(t, FourByteHeader, detectByteKind(0b11110111))

	assert.Equal(t, ContinuationByte, detectByteKind(0b10101010))
	assert.Equal(t, ContinuationByte, detectByteKind(0b10111111))
	assert.Equal(t, ContinuationByte, detectByteKind(0b10000000))
	assert.Equal(t, ContinuationByte, detectByteKind(0b10110101))
}

func Test_EndBytesMissing(t *testing.T) {
	assert.Equal(t, 0, detectEndBytesMissing([]byte{0x0A, 0x0B, 0x65, 0x11}))
	assert.Equal(t, 0, detectEndBytesMissing([]byte{0x0A, 0x0B, 0xC0, 0x8F}))
	assert.Equal(t, 0, detectEndBytesMissing([]byte{0x0A, 0xE0, 0x8F, 0x8F}))
	assert.Equal(t, 0, detectEndBytesMissing([]byte{0xF0, 0x8F, 0x8F, 0x8F}))

	assert.Equal(t, 1, detectEndBytesMissing([]byte{0x0A, 0x0B, 0x0A, 0xC0}))
	assert.Equal(t, 2, detectEndBytesMissing([]byte{0x0A, 0x0B, 0x0A, 0xE0}))
	assert.Equal(t, 3, detectEndBytesMissing([]byte{0x0A, 0x0B, 0x0A, 0xF0}))

	assert.Equal(t, 1, detectEndBytesMissing([]byte{0x0A, 0x0B, 0xE0, 0x8F}))
	assert.Equal(t, 2, detectEndBytesMissing([]byte{0x0A, 0x0B, 0xF0, 0x8F}))

	assert.Equal(t, 1, detectEndBytesMissing([]byte{0x00, 0xF0, 0x8F, 0x8F}))
}

func Test_DetectDanglingContinuationBytes(t *testing.T) {
	assert.Equal(t, 0, detectDanglingContinuationBytes([]byte{0x00, 0x00, 0x00, 0x00}))
	assert.Equal(t, 1, detectDanglingContinuationBytes([]byte{0x8F, 0x00, 0x00, 0x00}))
	assert.Equal(t, 2, detectDanglingContinuationBytes([]byte{0x8F, 0x8F, 0x00, 0x00}))
}
