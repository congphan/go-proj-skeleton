// +build unit

package strutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CleanEmpty(t *testing.T) {
	in := []string{"abc", "  ", "def", ""}

	assert.Equal(t, []string{"abc", "def"}, CleanEmpty(in))
}

func Test_Include(t *testing.T) {
	arrStrings := []string{"abc", "def"}

	assert.True(t, Include(arrStrings, "abc"))
	assert.True(t, Include(arrStrings, "def"))

	assert.False(t, Include(arrStrings, "abd"), `"abd" not belong to given array`)
}

func Test_Index(t *testing.T) {
	arrStrings := []string{"abc", "def"}

	assert.Equal(t, 0, Index(arrStrings, "abc"))
	assert.Equal(t, 1, Index(arrStrings, "def"))
	assert.Equal(t, -1, Index(arrStrings, "mk"), `"mk" not belong to given array`)
}
