package repo

import "go-prj-skeleton/app/domain/model"

type TransactionRepo interface {
	FindByUser(userID uint) ([]model.Transaction, error)
	FindByUserAccount(userID, accountID uint) ([]model.Transaction, error)
}
