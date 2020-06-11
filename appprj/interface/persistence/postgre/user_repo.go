package postgre

import (
	"github.com/go-pg/pg/v9"

	"go-prj-skeleton/appprj/domain/model"
	"go-prj-skeleton/appprj/pgutil"
)

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func toUser(u user) model.User {
	return model.User{
		ID:   u.ID,
		Name: u.Name,
	}
}

type userRepo struct {
}

func NewUserRepo() *userRepo {
	return &userRepo{}
}

func (repo userRepo) FindByID(id int) (model.User, error) {
	u := user{}

	_, err := pgutil.DB().QueryOne(&u, "SELECT * FROM users WHERE id=?", id)
	if err != nil {
		if err == pg.ErrNoRows {
			return model.User{}, model.ErrNotFound
		}

		return model.User{}, err
	}

	return toUser(u), nil
}
