package usecase

import (
	"encoding/json"
	"testing"

	"go-prj-skeleton/appprj/domain/model"

	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
)

func TestToTransactions(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		trans := []model.Transaction{
			{
				ID:              1,
				AccountID:       1,
				Amount:          decimal.NewFromFloat(10000),
				TransactionType: model.TransactionTypeDeposit,
				CreatedAt:       "2020-02-10 20:00:00 +0700",
			},
			{
				ID:              2,
				AccountID:       2,
				Amount:          decimal.NewFromFloat(20000),
				TransactionType: model.TransactionTypeWithdraw,
				CreatedAt:       "2020-02-12 20:00:00 +0700",
			},
		}

		accs := []model.Account{
			{
				ID:     1,
				UserID: 1,
				Name:   "Cong Phan",
				Bank:   "VCB",
			},
			{
				ID:     2,
				UserID: 1,
				Name:   "PHAN THANH CONG",
				Bank:   "ACB",
			},
		}

		out, err := toTransactions(trans, accs)
		assert.NoError(t, err)
		bytes, err := json.Marshal(out)
		assert.NoError(t, err)

		assert.JSONEq(t, `[
  {
    "ID": 1,
    "AccountID": 1,
    "Amount": "10000",
    "Bank": "VCB",
    "TransactionType": "deposit",
    "CreatedAt": "2020-02-10 20:00:00 +0700"
  },
  {
    "ID": 2,
    "AccountID": 2,
    "Amount": "20000",
    "Bank": "ACB",
    "TransactionType": "withdraw",
    "CreatedAt": "2020-02-12 20:00:00 +0700"
  }
]`, string(bytes))
	})
}
