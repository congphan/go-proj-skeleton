package repo

import "go-prj-skeleton/app/domain/model"

type AccountRepo interface {
	FindByUser(userID uint) ([]model.Account, error)
	FindByID(id uint) (model.Account, error)
}
