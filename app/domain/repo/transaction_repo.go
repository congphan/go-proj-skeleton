package repo

import "go-prj-skeleton/app/domain/model"

type TransactionRepo interface {
	FindByID(id int) (model.Transaction, error)
	FindByUser(userID int) ([]model.Transaction, error)
	FindByUserAccount(userID, accountID int) ([]model.Transaction, error)
	Create(*model.Transaction) error
	Update(*model.Transaction) error
	Delete(userID, tranID int) error
}
