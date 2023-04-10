package teststore

import (
	"time"

	"github.com/Aza-9798/costs-rest-api/internal/app/model"
	"github.com/Aza-9798/costs-rest-api/internal/app/store"
)

type TransactionRepository struct {
	store        *Store
	transactions map[int]*model.TransactionDB
}

func (r *TransactionRepository) Create(t *model.TransactionDB) error {
	if err := t.Validate(); err != nil {
		return err
	}

	t.ID = len(r.transactions)
	r.transactions[t.ID] = t
	return nil
}

func (r *TransactionRepository) Delete(t *model.TransactionJSON) error {
	if t1, ok := r.transactions[t.ID]; !ok {
		return store.ErrRecordNotFound
	} else {
		delete(r.transactions, t1.ID)
		return nil
	}
}

func (r *TransactionRepository) Save(t *model.TransactionDB) error {

	return nil
}

func (r *TransactionRepository) Find(id int) (*model.TransactionJSON, error) {
	t, ok := r.transactions[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return t.ToJSON(), nil
}

func (r *TransactionRepository) GetAllByAccount(accountID int) ([]*model.TransactionJSON, error) {
	res := make([]*model.TransactionJSON, 0)
	for _, t := range r.transactions {
		if t.Source.ID == accountID || t.Destination.ID == accountID {
			res = append(res, t.ToJSON())
		}
	}
	return res, nil
}

func (r *TransactionRepository) GetAllByAccountAndPeriod(userID int, DateStart, DateEnd time.Time) ([]*model.TransactionJSON, error) {
	return nil, nil
}

func (r *TransactionRepository) GetAllByUser(userID int) ([]*model.TransactionJSON, error) {
	return nil, nil
}

func (r *TransactionRepository) GetSummary(userID int, DateStart, DateEnd time.Time) (*model.Summary, error) {
	return nil, nil
}
