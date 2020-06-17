package model

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type TransactionType string

var (
	TransactionTypeWithdraw TransactionType = "withdraw"
	TransactionTypeDeposit  TransactionType = "deposit"

	ErrTransactionTypeInvalid = fmt.Errorf("invalid transaction type")
)

type Transaction struct {
	ID int

	AccountID int
	UserID    int

	Amount          decimal.Decimal
	TransactionType TransactionType
	CreatedAt       string
}

func ValidateTransactionType(t TransactionType) error {
	switch true {
	case t == TransactionTypeWithdraw:
		return nil

	case t == TransactionTypeDeposit:
		return nil

	default:
		return fmt.Errorf("%s: %w", t, ErrTransactionTypeInvalid)
	}
}

func NewTransaction(userID, accountID int, amount decimal.Decimal, t TransactionType) *Transaction {
	return &Transaction{
		UserID:          userID,
		AccountID:       accountID,
		Amount:          amount,
		TransactionType: t,
	}
}
