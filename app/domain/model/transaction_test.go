package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateTransactionType(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		assert.NoError(t, ValidateTransactionType(TransactionTypeDeposit))
		assert.NoError(t, ValidateTransactionType(TransactionTypeWithdraw))
	})

	t.Run("invalid", func(t *testing.T) {
		var err error
		err = ValidateTransactionType(TransactionType(""))
		errors.Is(err, ErrTransactionTypeInvalid)

		err = ValidateTransactionType(TransactionType("xyz"))
		errors.Is(err, ErrTransactionTypeInvalid)
	})
}
