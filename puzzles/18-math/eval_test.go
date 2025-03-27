package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	assert.Equal(t, float64(10), evalExpression("5 + 5"))
	assert.Equal(t, float64(55), evalExpression("5 * (6 + 5)"))
	assert.Equal(t, float64(-5), evalExpression("5 - 10"))
	assert.Equal(t, float64(66), evalExpression("((1 + 1) + 1) * ((4 - (15 - (66 / 2))) * 1)"))
}
