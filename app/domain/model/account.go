package model

type Account struct {
	ID int

	UserID int

	Name string
	Bank string
}

type Accounts []Account

func (s Accounts) ByID(id int) (Account, bool) {
	for _, acc := range s {
		if acc.ID == id {
			return acc, true
		}
	}

	return Account{}, false
}
