package repo

import "go-prj-skeleton/app/domain/model"

type TransactionRepo interface {
	FindByID(id uint) (model.Transaction, error)
	FindByUser(userID uint) ([]model.Transaction, error)
	FindByUserAccount(userID, accountID uint) ([]model.Transaction, error)
	Create(*model.Transaction) error
	Update(*model.Transaction) error
	Delete(userID, tranID uint) error
}
