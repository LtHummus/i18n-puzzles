package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseJapaneseNumber(t *testing.T) {
	t.Run("happy cases", func(t *testing.T) {
		tests := map[string]int64{
			"三百":                               300,
			"三百二十一":                         321,
			"四千":                               4_000,
			"五万":                               50_000,
			"九万九千九百九十九":                 99_999,
			"四十二万四十二":                     420_042,
			"九億八千七百六十五万四千三百二十一": 987_654_321,
		}

		for jaNum, val := range tests {
			computed, err := parseJapaneseNumber(jaNum)
			assert.NoError(t, err)
			assert.Equal(t, val, computed)
		}
	})

	t.Run("error cases", func(t *testing.T) {
		_, err := parseJapaneseNumber("三百二x十一")
		assert.Error(t, err)
	})
}
