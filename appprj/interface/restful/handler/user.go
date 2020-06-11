package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
	"goji.io/v3/pat"

	"go-prj-skeleton/appprj/domain/model"
	"go-prj-skeleton/appprj/jsonutil"
	"go-prj-skeleton/appprj/usecase"
)

type createTransaction struct {
	AccountID       int                   `json:"account_id"`
	Amount          decimal.Decimal       `json:"amount"`
	TransactionType model.TransactionType `json:"transaction_type"`
}

type UpdateTransaction struct {
	Amount decimal.Decimal `json:"amount"`
}

type transaction struct {
	ID              int                   `json:"id"`
	AccountID       int                   `json:"account_id"`
	Amount          decimal.Decimal       `json:"amount"`
	Bank            string                `json:"bank"`
	TransactionType model.TransactionType `json:"transaction_type"`
	CreatedAt       string                `json:"created_at"`
}

func toTransaction(t usecase.Transaction) transaction {
	return transaction{
		ID:              t.ID,
		AccountID:       t.AccountID,
		Amount:          t.Amount,
		Bank:            t.Bank,
		TransactionType: t.TransactionType,
		CreatedAt:       t.CreatedAt,
	}
}

func toTransactions(s []usecase.Transaction) []transaction {
	out := make([]transaction, len(s))

	for i := range s {
		out[i] = toTransaction(s[i])
	}

	return out
}

type userHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *userHandler {
	return &userHandler{
		userUsecase,
	}
}

type erroMessage struct {
	Error string `json:"error"`
}

func (h userHandler) FindTransactions(w http.ResponseWriter, r *http.Request) {
	strUserID := pat.Param(r, "user_id")
	userID, err := strconv.ParseInt(strUserID, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonutil.Marshal(erroMessage{err.Error()}))
		return
	}

	var accountID *int
	strAccountID := r.URL.Query().Get("account_id")
	if strAccountID != "" {
		accID, err := strconv.ParseInt(strAccountID, 10, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonutil.Marshal(erroMessage{err.Error()}))
			return
		}

		parsedID := int(accID)
		accountID = &parsedID
	}

	trans, err := h.userUsecase.FindTransactions(int(userID), accountID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonutil.Marshal(erroMessage{err.Error()}))
		return
	}

	bytes, err := json.Marshal(toTransactions(trans))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonutil.Marshal(erroMessage{err.Error()}))
		return
	}

	w.Write(bytes)
}

func (h userHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	strUserID := pat.Param(r, "user_id")
	pUserID, err := strconv.ParseInt(strUserID, 10, 32)
	userID := int(pUserID)
	if err != nil {
		w.Write(jsonutil.Marshal(erroMessage{err.Error()}))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payl := createTransaction{}
	if err := json.NewDecoder(r.Body).Decode(&payl); err != nil {
		w.Write(jsonutil.Marshal(erroMessage{err.Error()}))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdTran, err := h.userUsecase.CreateTransaction(userID, usecase.CreateTransaction{
		AccountID:       payl.AccountID,
		Amount:          payl.Amount,
		TransactionType: payl.TransactionType,
	})
	if err != nil {
		w.Write(jsonutil.Marshal(erroMessage{err.Error()}))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(toTransaction(*createdTran))
	if err != nil {
		w.Write(jsonutil.Marshal(erroMessage{err.Error()}))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
	w.WriteHeader(http.StatusCreated)
}

func (h userHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	strUserID := pat.Param(r, "user_id")
	userID, err := strconv.ParseInt(strUserID, 10, 32)
	if err != nil {
		w.Write(jsonutil.Marshal(erroMessage{err.Error()}))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	strTranID := pat.Param(r, "transaction_id")
	tranID, err := strconv.ParseInt(strTranID, 10, 32)
	if err != nil {
		w.Write(jsonutil.Marshal(erroMessage{err.Error()}))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payl := UpdateTransaction{}
	if err := json.NewDecoder(r.Body).Decode(&payl); err != nil {
		w.Write(jsonutil.Marshal(erroMessage{err.Error()}))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedTran, err := h.userUsecase.UpdateTransaction(int(userID), int(tranID), usecase.UpdateTransaction{
		Amount: payl.Amount,
	})
	if err != nil {
		w.Write(jsonutil.Marshal(erroMessage{err.Error()}))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(toTransaction(*updatedTran))
	if err != nil {
		w.Write(jsonutil.Marshal(erroMessage{err.Error()}))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (h userHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	strUserID := pat.Param(r, "user_id")
	userID, err := strconv.ParseInt(strUserID, 10, 32)
	if err != nil {
		w.Write(jsonutil.Marshal(erroMessage{err.Error()}))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	strTranID := pat.Param(r, "transaction_id")
	tranID, err := strconv.ParseInt(strTranID, 10, 32)
	if err != nil {
		w.Write(jsonutil.Marshal(erroMessage{err.Error()}))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.userUsecase.DeleteTransaction(int(userID), int(tranID))
	if err != nil {
		w.Write(jsonutil.Marshal(erroMessage{err.Error()}))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
