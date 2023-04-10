package teststore

import (
	"github.com/Aza-9798/costs-rest-api/internal/app/model"
	"github.com/Aza-9798/costs-rest-api/internal/app/store"
)

type AccountRepository struct {
	store    *Store
	accounts map[int]*model.Account
}

func (r *AccountRepository) Create(a *model.Account) error {
	if err := a.Validate(); err != nil {
		return err
	}

	a.ID = len(r.accounts)
	r.accounts[a.ID] = a
	return nil
}

func (r *AccountRepository) Save(a *model.Account) error {
	if err := a.Validate(); err != nil {
		return err
	}
	if _, ok := r.accounts[a.ID]; !ok {
		return store.ErrRecordNotFound
	}
	r.accounts[a.ID] = a
	return nil
}

func (r *AccountRepository) Find(id int) (*model.Account, error) {
	a, ok := r.accounts[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return a, nil
}

func (r *AccountRepository) Delete(id int) error {
	if _, ok := r.accounts[id]; !ok {
		return store.ErrRecordNotFound
	}
	delete(r.accounts, id)
	return nil
}

func (r *AccountRepository) GetAllByUser(userID int) ([]*model.Account, error) {
	res := make([]*model.Account, 0)
	for _, acc := range r.accounts {
		if acc.User == userID {
			res = append(res, acc)
		}
	}
	return res, nil
}

func (r *AccountRepository) IsAccountBelongsUser(accountID, userID int) (bool, error) {
	a, err := r.Find(accountID)
	if err != nil {
		return false, err
	}
	return a.User == userID, nil
}
