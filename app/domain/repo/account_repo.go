package repo

import "go-prj-skeleton/app/domain/model"

type AccountRepo interface {
	FindByUser(userID int) ([]model.Account, error)
	FindByID(id int) (model.Account, error)
}
