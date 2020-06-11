package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"go-prj-skeleton/appprj/domain/model"
	"go-prj-skeleton/appprj/domain/repo/mock"

	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
)

func TestUserUsecase_FindTransactions(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		userRepo := &mock.FakeUserRepo{
			FindByIDHook: func(userID int) (model.User, error) {
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
			FindByUserHook: func(userID int) ([]model.Transaction, error) {
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
			FindByUserAccountHook: func(userID int, accountID int) ([]model.Transaction, error) {
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
			FindByUserHook: func(userID int) ([]model.Account, error) {
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

			accountID := int(1)
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
			accountID := int(2)
			trans, err := uc.FindTransactions(1, &accountID)
			assert.NoError(t, err)
			assert.Equal(t, 0, len(trans))
		})
	})

	t.Run("invalid", func(t *testing.T) {
		t.Parallel()

		userRepo := &mock.FakeUserRepo{
			FindByIDHook: func(userID int) (model.User, error) {
				if userID == 3 {
					return model.User{3, "John"}, nil
				}

				return model.User{}, fmt.Errorf("user id:%v %w", userID, model.ErrNotFound)
			},
		}

		transRepo := &mock.FakeTransactionRepo{
			FindByUserHook: func(userID int) ([]model.Transaction, error) {
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
			FindByUserAccountHook: func(userID int, accountID int) ([]model.Transaction, error) {
				return nil, nil
			},
		}

		accountRepo := &mock.FakeAccountRepo{
			FindByUserHook: func(userID int) ([]model.Account, error) {
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
			FindByIDHook: func(userID int) (model.User, error) {
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
			FindByUserHook: func(userID int) ([]model.Transaction, error) {
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
			FindByUserAccountHook: func(userID int, accountID int) ([]model.Transaction, error) {
				if userID == 1 && accountID == 1 {
					return nil, fmt.Errorf("find transactions by user and account got internal error")
				}

				return nil, nil
			},
		}

		accountRepo := &mock.FakeAccountRepo{
			FindByUserHook: func(userID int) ([]model.Account, error) {
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
			accountID := int(1)
			_, err := uc.FindTransactions(1, &accountID)
			assert.EqualError(t, err, "find transactions by user and account got internal error")
		})

		t.Run("find accounts by user", func(t *testing.T) {
			_, err := uc.FindTransactions(3, nil)
			assert.EqualError(t, err, "find accounts by user got internal error")
		})
	})
}

func TestUserUsecase_CreateTransaction(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		userRepo := &mock.FakeUserRepo{
			FindByIDHook: func(userID int) (model.User, error) {
				if userID == 1 {
					return model.User{
						1,
						"Alice",
					}, nil
				}

				return model.User{}, model.ErrNotFound
			},
		}

		accountRepo := &mock.FakeAccountRepo{
			FindByIDHook: func(accountID int) (model.Account, error) {
				if accountID == 1 {
					return model.Account{
						ID:     1,
						UserID: 1,
						Name:   "Alice A",
						Bank:   "VCB",
					}, nil
				}

				return model.Account{}, nil
			},
		}

		tranRepo := &mock.FakeTransactionRepo{
			CreateHook: func(t *model.Transaction) error {
				t.ID = 123
				t.CreatedAt = "2020-02-10 20:10:00 +0700"

				return nil
			},
		}

		uc := NewUserUsecase(userRepo, accountRepo, tranRepo)
		createdTran, err := uc.CreateTransaction(1, CreateTransaction{
			AccountID:       1,
			Amount:          decimal.NewFromInt(1000),
			TransactionType: model.TransactionTypeDeposit,
		})
		assert.NoError(t, err)

		bytes, err := json.Marshal(createdTran)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
  "ID": 123,
  "AccountID": 1,
  "Amount": "1000",
  "Bank": "VCB",
  "TransactionType": "deposit",
  "CreatedAt": "2020-02-10 20:10:00 +0700"
}`, string(bytes))

	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()

		t.Run("invalid transaction type", func(t *testing.T) {
			tran := CreateTransaction{
				AccountID:       1,
				Amount:          decimal.NewFromInt(1000),
				TransactionType: "TTT",
			}

			uc := NewUserUsecase(nil, nil, nil)
			_, err := uc.CreateTransaction(1, tran)
			assert.EqualError(t, err, "TTT: invalid transaction type")
		})

		t.Run("invalid ammount", func(t *testing.T) {
			tran := CreateTransaction{
				AccountID:       1,
				Amount:          decimal.NewFromInt(0),
				TransactionType: model.TransactionTypeDeposit,
			}

			uc := NewUserUsecase(nil, nil, nil)
			_, err := uc.CreateTransaction(1, tran)
			assert.EqualError(t, err, "amount[0]: invalid")
		})

		t.Run("find user by id fail", func(t *testing.T) {
			tran := CreateTransaction{
				AccountID:       1,
				Amount:          decimal.NewFromInt(1000),
				TransactionType: model.TransactionTypeDeposit,
			}

			userRepo := &mock.FakeUserRepo{
				FindByIDHook: func(int) (model.User, error) {
					return model.User{}, model.ErrNotFound
				},
			}

			uc := NewUserUsecase(userRepo, nil, nil)
			_, err := uc.CreateTransaction(1, tran)
			assert.True(t, errors.Is(err, model.ErrNotFound))
			assert.EqualError(t, err, "not found")
		})

		t.Run("find account by id fail", func(t *testing.T) {
			tran := CreateTransaction{
				AccountID:       1,
				Amount:          decimal.NewFromInt(1000),
				TransactionType: model.TransactionTypeDeposit,
			}

			userRepo := &mock.FakeUserRepo{
				FindByIDHook: func(userID int) (model.User, error) {
					if userID == 1 {
						return model.User{1, "Alice"}, nil
					}

					return model.User{}, model.ErrNotFound
				},
			}

			accountRepo := &mock.FakeAccountRepo{
				FindByIDHook: func(int) (model.Account, error) {
					return model.Account{}, model.ErrNotFound
				},
			}

			uc := NewUserUsecase(userRepo, accountRepo, nil)
			_, err := uc.CreateTransaction(1, tran)
			assert.True(t, errors.Is(err, model.ErrInvalid))
			assert.EqualError(t, err, "account[1] invalid")
		})

		t.Run("account not belong to user", func(t *testing.T) {
			tran := CreateTransaction{
				AccountID:       1,
				Amount:          decimal.NewFromInt(1000),
				TransactionType: model.TransactionTypeDeposit,
			}

			userRepo := &mock.FakeUserRepo{
				FindByIDHook: func(userID int) (model.User, error) {
					if userID == 1 {
						return model.User{1, "Alice"}, nil
					}

					return model.User{}, model.ErrNotFound
				},
			}

			accountRepo := &mock.FakeAccountRepo{
				FindByIDHook: func(accountID int) (model.Account, error) {
					if accountID == 1 {
						return model.Account{
							ID:     1,
							UserID: 2,
						}, nil
					}

					return model.Account{}, model.ErrNotFound
				},
			}

			uc := NewUserUsecase(userRepo, accountRepo, nil)
			_, err := uc.CreateTransaction(1, tran)
			assert.True(t, errors.Is(err, model.ErrInvalid))
			assert.EqualError(t, err, "account[1] invalid")
		})

		t.Run("something wrong when persisting transaction", func(t *testing.T) {
			tran := CreateTransaction{
				AccountID:       1,
				Amount:          decimal.NewFromInt(1000),
				TransactionType: model.TransactionTypeDeposit,
			}

			userRepo := &mock.FakeUserRepo{
				FindByIDHook: func(userID int) (model.User, error) {
					if userID == 1 {
						return model.User{1, "Alice"}, nil
					}

					return model.User{}, model.ErrNotFound
				},
			}

			accountRepo := &mock.FakeAccountRepo{
				FindByIDHook: func(accountID int) (model.Account, error) {
					if accountID == 1 {
						return model.Account{
							ID:     1,
							UserID: 1,
						}, nil
					}

					return model.Account{}, model.ErrNotFound
				},
			}

			tranRepo := &mock.FakeTransactionRepo{
				CreateHook: func(*model.Transaction) error {
					return fmt.Errorf("internal error")
				},
			}

			uc := NewUserUsecase(userRepo, accountRepo, tranRepo)
			_, err := uc.CreateTransaction(1, tran)
			assert.EqualError(t, err, "persit transaction: internal error")
		})
	})
}

func TestUserUsecase_UpdateTransaction(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		userRepo := &mock.FakeUserRepo{
			FindByIDHook: func(userID int) (model.User, error) {
				if userID == 1 {
					return model.User{
						1, "Alice",
					}, nil
				}

				return model.User{}, model.ErrNotFound
			},
		}

		tranRepo := &mock.FakeTransactionRepo{
			FindByIDHook: func(tranID int) (model.Transaction, error) {
				if tranID == 2 {
					return model.Transaction{
						ID:              2,
						AccountID:       3,
						Amount:          decimal.NewFromInt(1000),
						TransactionType: model.TransactionTypeDeposit,
						CreatedAt:       "2020-02-10 20:10:00 +0700",
					}, nil
				}

				return model.Transaction{}, model.ErrNotFound
			},
			UpdateHook: func(t *model.Transaction) error {
				return nil
			},
		}

		accountRepo := &mock.FakeAccountRepo{
			FindByUserHook: func(userID int) ([]model.Account, error) {
				if userID == 1 {
					return []model.Account{
						{
							ID:     3,
							UserID: 1,
							Name:   "PHAN THANH CONG",
							Bank:   "ACB",
						},
					}, nil
				}

				return nil, nil
			},
		}

		uc := NewUserUsecase(userRepo, accountRepo, tranRepo)

		tran, err := uc.UpdateTransaction(1, 2, UpdateTransaction{decimal.NewFromInt(2000)})
		assert.NoError(t, err)

		bytes, err := json.Marshal(tran)
		assert.NoError(t, err)
		assert.JSONEq(t, `{
  "ID": 2,
  "AccountID": 3,
  "Amount": "2000",
  "Bank": "ACB",
  "TransactionType": "deposit",
  "CreatedAt": "2020-02-10 20:10:00 +0700"
}`, string(bytes))
	})

	t.Run("fail", func(t *testing.T) {
		t.Run("zero amount", func(t *testing.T) {
			uc := NewUserUsecase(nil, nil, nil)
			_, err := uc.UpdateTransaction(1, 2, UpdateTransaction{decimal.NewFromInt(0)})
			assert.True(t, errors.Is(err, model.ErrInvalid))
			assert.EqualError(t, err, "amount[0]: invalid")
		})

		t.Run("find user fail", func(t *testing.T) {
			userRepo := &mock.FakeUserRepo{
				FindByIDHook: func(userID int) (model.User, error) {
					return model.User{}, model.ErrNotFound
				},
			}

			uc := NewUserUsecase(userRepo, nil, nil)

			_, err := uc.UpdateTransaction(1, 2, UpdateTransaction{decimal.NewFromInt(2000)})
			assert.True(t, errors.Is(err, model.ErrNotFound))
			assert.EqualError(t, err, "find user[1] not found")
		})

		t.Run("find transaction fail", func(t *testing.T) {
			userRepo := &mock.FakeUserRepo{
				FindByIDHook: func(userID int) (model.User, error) {
					if userID == 1 {
						return model.User{
							1, "Alice",
						}, nil
					}

					return model.User{}, model.ErrNotFound
				},
			}

			tranRepo := &mock.FakeTransactionRepo{
				FindByIDHook: func(tranID int) (model.Transaction, error) {
					return model.Transaction{}, model.ErrNotFound
				},
			}

			uc := NewUserUsecase(userRepo, nil, tranRepo)

			_, err := uc.UpdateTransaction(1, 2, UpdateTransaction{decimal.NewFromInt(2000)})
			assert.True(t, errors.Is(err, model.ErrNotFound))
			assert.EqualError(t, err, "find transaction[2] not found")
		})

		t.Run("find account fail", func(t *testing.T) {
			userRepo := &mock.FakeUserRepo{
				FindByIDHook: func(userID int) (model.User, error) {
					if userID == 1 {
						return model.User{
							1, "Alice",
						}, nil
					}

					return model.User{}, model.ErrNotFound
				},
			}

			tranRepo := &mock.FakeTransactionRepo{
				FindByIDHook: func(tranID int) (model.Transaction, error) {
					if tranID == 2 {
						return model.Transaction{
							ID:              2,
							AccountID:       3,
							Amount:          decimal.NewFromInt(1000),
							TransactionType: model.TransactionTypeDeposit,
							CreatedAt:       "2020-02-10 20:10:00 +0700",
						}, nil
					}

					return model.Transaction{}, model.ErrNotFound
				},
				UpdateHook: func(t *model.Transaction) error {
					return nil
				},
			}

			accountRepo := &mock.FakeAccountRepo{
				FindByUserHook: func(userID int) ([]model.Account, error) {
					return nil, fmt.Errorf("internal error")
				},
			}

			uc := NewUserUsecase(userRepo, accountRepo, tranRepo)

			_, err := uc.UpdateTransaction(1, 2, UpdateTransaction{decimal.NewFromInt(2000)})
			assert.EqualError(t, err, "internal error")
		})

		t.Run("invalid transaction", func(t *testing.T) {
			userRepo := &mock.FakeUserRepo{
				FindByIDHook: func(userID int) (model.User, error) {
					if userID == 1 {
						return model.User{
							1, "Alice",
						}, nil
					}

					return model.User{}, model.ErrNotFound
				},
			}

			tranRepo := &mock.FakeTransactionRepo{
				FindByIDHook: func(tranID int) (model.Transaction, error) {
					if tranID == 2 {
						return model.Transaction{
							ID:              2,
							AccountID:       4,
							Amount:          decimal.NewFromInt(1000),
							TransactionType: model.TransactionTypeDeposit,
							CreatedAt:       "2020-02-10 20:10:00 +0700",
						}, nil
					}

					return model.Transaction{}, model.ErrNotFound
				},
				UpdateHook: func(t *model.Transaction) error {
					return nil
				},
			}

			accountRepo := &mock.FakeAccountRepo{
				FindByUserHook: func(userID int) ([]model.Account, error) {
					if userID == 1 {
						return []model.Account{
							{
								ID:     3,
								UserID: 1,
								Name:   "PHAN THANH CONG",
								Bank:   "ACB",
							},
						}, nil
					}

					return nil, nil
				},
			}

			uc := NewUserUsecase(userRepo, accountRepo, tranRepo)

			_, err := uc.UpdateTransaction(1, 2, UpdateTransaction{decimal.NewFromInt(2000)})
			assert.True(t, errors.Is(err, model.ErrInvalid))
			assert.EqualError(t, err, "transaction[2] invalid")
		})

		t.Run("something wrong when persisting transaction", func(t *testing.T) {
			userRepo := &mock.FakeUserRepo{
				FindByIDHook: func(userID int) (model.User, error) {
					if userID == 1 {
						return model.User{
							1, "Alice",
						}, nil
					}

					return model.User{}, model.ErrNotFound
				},
			}

			tranRepo := &mock.FakeTransactionRepo{
				FindByIDHook: func(tranID int) (model.Transaction, error) {
					if tranID == 2 {
						return model.Transaction{
							ID:              2,
							AccountID:       3,
							Amount:          decimal.NewFromInt(1000),
							TransactionType: model.TransactionTypeDeposit,
							CreatedAt:       "2020-02-10 20:10:00 +0700",
						}, nil
					}

					return model.Transaction{}, model.ErrNotFound
				},
				UpdateHook: func(t *model.Transaction) error {
					return fmt.Errorf("internal error")
				},
			}

			accountRepo := &mock.FakeAccountRepo{
				FindByUserHook: func(userID int) ([]model.Account, error) {
					if userID == 1 {
						return []model.Account{
							{
								ID:     3,
								UserID: 1,
								Name:   "PHAN THANH CONG",
								Bank:   "ACB",
							},
						}, nil
					}

					return nil, nil
				},
			}

			uc := NewUserUsecase(userRepo, accountRepo, tranRepo)

			_, err := uc.UpdateTransaction(1, 2, UpdateTransaction{decimal.NewFromInt(2000)})
			assert.EqualError(t, err, "update transaction[2] internal error")
		})
	})
}
