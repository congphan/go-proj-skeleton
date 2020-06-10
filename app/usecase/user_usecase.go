package usecase

import (
	"go-prj-skeleton/app/domain/model"
	"go-prj-skeleton/app/domain/repo"
)

type UserUsecase interface {
	FindTransactions(userID uint, accountID *uint) ([]Transaction, error)
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
