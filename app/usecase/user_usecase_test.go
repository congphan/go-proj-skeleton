package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"go-prj-skeleton/app/domain/model"
	"go-prj-skeleton/app/domain/repo/mock"

	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
)

func TestUserUsecase_FindTransactions(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		userRepo := &mock.FakeUserRepo{
			FindByIDHook: func(userID uint) (model.User, error) {
				if userID == 1 {
					return model.User{1, "Cong Phan"}, nil
				}

				if userID == 2 {
					return model.User{2, "Alice"}, nil
				}

				return model.User{}, fmt.Errorf("user id:%v %w", userID, model.ErrNotFound)
			},
		}

		transRepo := &mock.FakeTransactionRepo{
			FindByUserHook: func(userID uint) ([]model.Transaction, error) {
				if userID == 1 {
					return []model.Transaction{
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
					}, nil
				}

				return nil, nil
			},
			FindByUserAccountHook: func(userID uint, accountID uint) ([]model.Transaction, error) {
				if userID == 1 && accountID == 1 {
					return []model.Transaction{
						{
							ID:              1,
							AccountID:       1,
							Amount:          decimal.NewFromFloat(10000),
							TransactionType: model.TransactionTypeDeposit,
							CreatedAt:       "2020-02-10 20:00:00 +0700",
						},
					}, nil
				}

				return nil, nil
			},
		}

		accountRepo := &mock.FakeAccountRepo{
			FindByUserHook: func(userID uint) ([]model.Account, error) {
				if userID == 1 {
					return []model.Account{
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
					}, nil
				}

				return nil, nil
			},
		}

		uc := NewUserUsecase(userRepo, accountRepo, transRepo)

		t.Run("valid user & empty account id", func(t *testing.T) {
			t.Parallel()

			trans, err := uc.FindTransactions(1, nil)
			assert.NoError(t, err)

			bytes, err := json.Marshal(trans)
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

		t.Run("valid user & existed account id", func(t *testing.T) {
			t.Parallel()

			accountID := uint(1)
			trans, err := uc.FindTransactions(1, &accountID)
			assert.NoError(t, err)

			bytes, err := json.Marshal(trans)
			assert.NoError(t, err)
			assert.JSONEq(t, `[
  {
    "ID": 1,
    "AccountID": 1,
    "Amount": "10000",
    "Bank": "VCB",
    "TransactionType": "deposit",
    "CreatedAt": "2020-02-10 20:00:00 +0700"
  }]`, string(bytes))
		})

		t.Run("valid user with no transactions", func(t *testing.T) {
			trans, err := uc.FindTransactions(2, nil)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(trans))
		})

		t.Run("valid user & account_id has no transaction", func(t *testing.T) {
			accountID := uint(2)
			trans, err := uc.FindTransactions(1, &accountID)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(trans))
		})
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		userRepo := &mock.FakeUserRepo{
			FindByIDHook: func(userID uint) (model.User, error) {
				if userID == 3 {
					return model.User{3, "John"}, nil
				}

				return model.User{}, fmt.Errorf("user id:%v %w", userID, model.ErrNotFound)
			},
		}

		transRepo := &mock.FakeTransactionRepo{
			FindByUserHook: func(userID uint) ([]model.Transaction, error) {
				if userID == 3 {
					return []model.Transaction{
						{
							ID:              3,
							AccountID:       4,
							Amount:          decimal.NewFromFloat(10000),
							TransactionType: model.TransactionTypeDeposit,
							CreatedAt:       "2020-02-10 20:00:00 +0700",
						},
					}, nil
				}

				return nil, nil
			},
			FindByUserAccountHook: func(userID uint, accountID uint) ([]model.Transaction, error) {
				return nil, nil
			},
		}

		accountRepo := &mock.FakeAccountRepo{
			FindByUserHook: func(userID uint) ([]model.Account, error) {
				return nil, nil
			},
		}

		uc := NewUserUsecase(userRepo, accountRepo, transRepo)

		t.Run("user has transaction but contains invalid account id", func(t *testing.T) {
			_, err := uc.FindTransactions(3, nil)
			if assert.Error(t, err) {
				assert.EqualError(t, err, "account[4] not found")
			}
		})

		t.Run("user not found", func(t *testing.T) {
			_, err := uc.FindTransactions(4, nil)
			if assert.Error(t, err) {
				assert.True(t, errors.Is(err, model.ErrNotFound))
				assert.EqualError(t, err, "user id:4 not found")
			}
		})
	})

	t.Run("unexpected error", func(t *testing.T) {
		userRepo := &mock.FakeUserRepo{
			FindByIDHook: func(userID uint) (model.User, error) {
				if userID == 1 {
					return model.User{1, "Cong Phan"}, nil
				}

				if userID == 2 {
					return model.User{2, "Alice"}, nil
				}

				if userID == 3 {
					return model.User{3, "John"}, nil
				}

				return model.User{}, fmt.Errorf("user id:%v %w", userID, model.ErrNotFound)
			},
		}

		transRepo := &mock.FakeTransactionRepo{
			FindByUserHook: func(userID uint) ([]model.Transaction, error) {
				if userID == 1 {
					return nil, fmt.Errorf("find transactions by user got internal error")
				}

				if userID == 3 {
					return []model.Transaction{
						{
							ID:              3,
							AccountID:       4,
							Amount:          decimal.NewFromFloat(10000),
							TransactionType: model.TransactionTypeDeposit,
							CreatedAt:       "2020-02-10 20:00:00 +0700",
						},
					}, nil
				}

				return nil, nil
			},
			FindByUserAccountHook: func(userID uint, accountID uint) ([]model.Transaction, error) {
				if userID == 1 && accountID == 1 {
					return nil, fmt.Errorf("find transactions by user and account got internal error")
				}

				return nil, nil
			},
		}

		accountRepo := &mock.FakeAccountRepo{
			FindByUserHook: func(userID uint) ([]model.Account, error) {
				if userID == 3 {
					return nil, fmt.Errorf("find accounts by user got internal error")
				}

				return nil, nil
			},
		}

		uc := NewUserUsecase(userRepo, accountRepo, transRepo)

		t.Run("find transaction by user", func(t *testing.T) {
			_, err := uc.FindTransactions(1, nil)
			assert.EqualError(t, err, "find transactions by user got internal error")
		})

		t.Run("find transaction by user and account", func(t *testing.T) {
			accountID := uint(1)
			_, err := uc.FindTransactions(1, &accountID)
			assert.EqualError(t, err, "find transactions by user and account got internal error")
		})

		t.Run("find accounts by user", func(t *testing.T) {
			_, err := uc.FindTransactions(3, nil)
			assert.EqualError(t, err, "find accounts by user got internal error")
		})
	})
}
