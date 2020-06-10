package model

type Account struct {
	ID uint

	UserID uint

	Name string
	Bank string
}

type Accounts []Account

func (s Accounts) ByID(id uint) (Account, bool) {
	for _, acc := range s {
		if acc.ID == id {
			return acc, true
		}
	}

	return Account{}, false
}
