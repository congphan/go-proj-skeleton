package postgre

import (
	"github.com/go-pg/pg/v9"

	"go-prj-skeleton/appprj/domain/model"
	"go-prj-skeleton/appprj/pgutil"
)

type account struct {
	ID int `json:"id"`

	UserID int `json:"user_id"`

	Name string `json:"name"`
	Bank string `json:"bank"`
}

func toAccount(acc account) model.Account {
	return model.Account{
		ID:     acc.ID,
		UserID: acc.UserID,
		Name:   acc.Name,
		Bank:   acc.Bank,
	}
}

type accountRepo struct {
}

func NewAccountRepo() *accountRepo {
	return &accountRepo{}
}

func (repo *accountRepo) FindByUser(userID int) ([]model.Account, error) {
	accs := []account{}

	_, err := pgutil.DB().Query(&accs, "SELECT * FROM accounts WHERE user_id=?", userID)
	if err != nil {
		return nil, err
	}

	out := make([]model.Account, len(accs))
	for i := range accs {
		out[i] = toAccount(accs[i])
	}

	return out, nil
}

func (repo *accountRepo) FindByID(id int) (model.Account, error) {
	acc := account{}

	_, err := pgutil.DB().QueryOne(&acc, "SELECT * FROM accounts WHERE id=?", id)
	if err != nil {
		if err == pg.ErrNoRows {
			return model.Account{}, model.ErrNotFound
		}

		return model.Account{}, err
	}

	return toAccount(acc), nil
}
