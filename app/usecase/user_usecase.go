package usecase

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"

	"go-prj-skeleton/app/domain/model"
	"go-prj-skeleton/app/domain/repo"
)

type UserUsecase interface {
	FindTransactions(userID uint, accountID *uint) ([]Transaction, error)
	CreateTransaction(userID uint, t CreateTransaction) (*Transaction, error)
	UpdateTransaction(userID, tranID uint, t UpdateTransaction) (*Transaction, error)
	DeleteTransaction(userID, tranID uint) error
}

type userUsecase struct {
	userRepo    repo.UserRepo
	accountRepo repo.AccountRepo
	transRepo   repo.TransactionRepo
}

func NewUserUsecase(userRepo repo.UserRepo, accountRepo repo.AccountRepo, transRepo repo.TransactionRepo) *userUsecase {
	return &userUsecase{
		userRepo,
		accountRepo,
		transRepo,
	}
}

func (u *userUsecase) FindTransactions(userID uint, accountID *uint) ([]Transaction, error) {
	_, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	var trans []model.Transaction
	if accountID == nil {
		trans, err = u.transRepo.FindByUser(userID)
		if err != nil {
			return nil, err
		}
	}

	if accountID != nil {
		trans, err = u.transRepo.FindByUserAccount(userID, *accountID)
		if err != nil {
			return nil, err
		}
	}

	if len(trans) == 0 {
		return []Transaction{}, nil
	}

	accounts, err := u.accountRepo.FindByUser(userID)
	if err != nil {
		return nil, err
	}

	return toTransactions(trans, accounts)
}

func (u *userUsecase) CreateTransaction(userID uint, t CreateTransaction) (*Transaction, error) {
	if err := model.ValidateTransactionType(t.TransactionType); err != nil {
		return nil, err
	}

	zero := decimal.NewFromInt(0)
	if t.Amount.LessThanOrEqual(zero) {
		return nil, fmt.Errorf("amount[%v]: %w", t.Amount.String(), model.ErrInvalid)
	}

	_, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	acc, err := u.accountRepo.FindByID(t.AccountID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, fmt.Errorf("account[%v] %w", t.AccountID, model.ErrInvalid)
		}

		return nil, err
	}

	if acc.UserID != userID {
		return nil, fmt.Errorf("account[%v] %w", t.AccountID, model.ErrInvalid)
	}

	tran := model.NewTransaction(t.AccountID, t.Amount, t.TransactionType)
	if err := u.transRepo.Create(tran); err != nil {
		return nil, fmt.Errorf("persit transaction: %w", err)
	}

	return &Transaction{
		ID:              tran.ID,
		AccountID:       tran.AccountID,
		Amount:          tran.Amount,
		Bank:            acc.Bank,
		TransactionType: tran.TransactionType,
		CreatedAt:       tran.CreatedAt,
	}, nil
}

func (u userUsecase) UpdateTransaction(userID, tranID uint, t UpdateTransaction) (*Transaction, error) {
	zero := decimal.NewFromInt(0)
	if t.Amount.LessThanOrEqual(zero) {
		return nil, fmt.Errorf("amount[%v]: %w", t.Amount.String(), model.ErrInvalid)
	}

	_, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("find user[%v] %w", userID, err)
	}

	tran, err := u.transRepo.FindByID(tranID)
	if err != nil {
		return nil, fmt.Errorf("find transaction[%v] %w", tranID, err)
	}

	accs, err := u.accountRepo.FindByUser(userID)
	if err != nil {
		return nil, err
	}

	acc, ok := model.Accounts(accs).ByID(tran.AccountID)
	if !ok {
		return nil, fmt.Errorf("transaction[%v] %w", tranID, model.ErrInvalid)
	}

	tran.Amount = t.Amount
	if err := u.transRepo.Update(&tran); err != nil {
		return nil, fmt.Errorf("update transaction[%v] %w", tranID, err)
	}

	return &Transaction{
		ID:              tran.ID,
		AccountID:       tran.AccountID,
		Amount:          tran.Amount,
		Bank:            acc.Bank,
		TransactionType: tran.TransactionType,
		CreatedAt:       tran.CreatedAt,
	}, nil
}

func (u *userUsecase) DeleteTransaction(userID, tranID uint) error {
	return u.transRepo.Delete(userID, tranID)
}
