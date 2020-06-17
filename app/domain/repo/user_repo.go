package repo

import "go-prj-skeleton/app/domain/model"

type UserRepo interface {
	FindByID(id int) (model.User, error)
}
