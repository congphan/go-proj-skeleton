package postgre

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/shopspring/decimal"

	"go-prj-skeleton/app/domain/model"
	"go-prj-skeleton/app/pgutil"
)

type transaction struct {
	ID int `json:"id"`

	UserID    int `json:"user_id"`
	AccountID int `json:"account_id"`

	Amount          decimal.Decimal       `json:"amount"`
	TransactionType model.TransactionType `json:"transaction_type"`
	CreatedAt       string                `json:"created_at"`
}

func toTransaction(t transaction) model.Transaction {
	return model.Transaction{
		ID:              t.ID,
		UserID:          t.UserID,
		AccountID:       t.AccountID,
		Amount:          t.Amount,
		TransactionType: t.TransactionType,
		CreatedAt:       t.CreatedAt,
	}
}

type transactionRepo struct {
}

func NewTransactionRepo() *transactionRepo {
	return &transactionRepo{}
}

func (repo transactionRepo) FindByID(id int) (model.Transaction, error) {
	tran := transaction{}

	_, err := pgutil.DB().QueryOne(&tran, "SELECT * FROM transactions WHERE id=?", id)
	if err != nil {
		if err == pg.ErrNoRows {
			return model.Transaction{}, model.ErrNotFound
		}

		return model.Transaction{}, err
	}

	return toTransaction(tran), nil
}

func (repo transactionRepo) FindByUser(userID int) ([]model.Transaction, error) {
	trans := []transaction{}

	_, err := pgutil.DB().Query(&trans, "SELECT * FROM transactions WHERE user_id=?", userID)
	if err != nil {
		return nil, err
	}

	out := make([]model.Transaction, len(trans))
	for i := range trans {
		out[i] = toTransaction(trans[i])
	}

	return out, nil
}

func (repo transactionRepo) FindByUserAccount(userID, accountID int) ([]model.Transaction, error) {
	trans := []transaction{}

	_, err := pgutil.DB().Query(&trans, "SELECT t.* FROM transactions t INNER JOIN accounts a ON t.account_id = a.id INNER JOIN users u ON a.user_id = u.id WHERE u.id=? AND a.id=?", userID, accountID)
	if err != nil {
		return nil, err
	}

	out := make([]model.Transaction, len(trans))
	for i := range trans {
		out[i] = toTransaction(trans[i])
	}

	return out, nil
}

func (repo transactionRepo) Create(t *model.Transaction) error {
	tran := transaction{
		ID:              int(time.Now().Unix()),
		AccountID:       t.AccountID,
		UserID:          t.UserID,
		Amount:          t.Amount,
		TransactionType: t.TransactionType,
	}

	now := time.Now().UTC().String()
	tran.CreatedAt = now
	if err := pgutil.DB().Insert(&tran); err != nil {
		return fmt.Errorf("exec Insert fail: %v", err)
	}

	t.CreatedAt = tran.CreatedAt
	t.ID = tran.ID

	return nil
}

func (repo transactionRepo) Update(t *model.Transaction) error {
	db := pgutil.DB()

	_, err := db.Model(&transaction{}).Set("amount=?", t.Amount).
		Where("id=?", t.ID).Update()
	if err != nil {
		return fmt.Errorf("update transaction fail: %v", err)
	}

	return nil
}

func (repo transactionRepo) Delete(userID, tranID int) error {
	tran, err := repo.FindByID(tranID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil
		}

		return err
	}

	if tran.UserID != userID {
		return nil
	}

	if err := pgHelper.delete(&transaction{ID: tran.ID}); err != nil {
		return err
	}

	return nil
}
