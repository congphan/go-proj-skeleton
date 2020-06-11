package repo

import "go-prj-skeleton/appprj/domain/model"

type UserRepo interface {
	FindByID(id int) (model.User, error)
}
