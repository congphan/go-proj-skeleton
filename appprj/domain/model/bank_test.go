package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateBank(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		// "VCB", "ACB", "VIB"
		assert.NoError(t, ValidateBank("VCB"))
		assert.NoError(t, ValidateBank("ACB"))
		assert.NoError(t, ValidateBank("VIB"))
	})

	t.Run("invalid", func(t *testing.T) {
		var err error

		err = ValidateBank("")
		errors.Is(err, ErrInvalidBank)
		assert.EqualError(t, err, ": invalid bank")

		err = ValidateBank("abc")
		errors.Is(err, ErrInvalidBank)
		assert.EqualError(t, err, "abc: invalid bank")
	})
}
