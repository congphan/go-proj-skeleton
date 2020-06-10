package usecase

import (
	"fmt"
	"go-prj-skeleton/app/domain/model"

	"github.com/shopspring/decimal"
)

type CreateTransaction struct {
	AccountID       uint
	Amount          decimal.Decimal
	TransactionType model.TransactionType
}

type Transaction struct {
	ID              uint
	AccountID       uint
	Amount          decimal.Decimal
	Bank            string
	TransactionType model.TransactionType
	CreatedAt       string
}

func toTransactions(trans []model.Transaction, accounts model.Accounts) ([]Transaction, error) {
	out := make([]Transaction, len(trans))

	for i := range trans {
		acc, ok := accounts.ByID(trans[i].AccountID)
		if !ok {
			return nil, fmt.Errorf("account[%v] %w", trans[i].AccountID, model.ErrNotFound)
		}

		t := Transaction{
			ID:              trans[i].ID,
			AccountID:       trans[i].AccountID,
			Amount:          trans[i].Amount,
			Bank:            acc.Bank,
			TransactionType: trans[i].TransactionType,
			CreatedAt:       trans[i].CreatedAt,
		}

		out[i] = t
	}

	return out, nil
}
