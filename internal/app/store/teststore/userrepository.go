package teststore

import (
	"first-rest-api/internal/app/model"
	"first-rest-api/internal/app/store"
)

type UserRepository struct {
	store *Store
	users map[int]*model.User
}

func (ur *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	if u1, _ := ur.FindByEmail(u.Email); u1 != nil {
		return store.ErrUserAlreadyExists
	}
	u.ID = len(ur.users)
	ur.users[u.ID] = u
	return nil
}

func (ur *UserRepository) Find(id int) (*model.User, error) {
	u, ok := ur.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}
	return u, nil
}

func (ur *UserRepository) FindByEmail(email string) (*model.User, error) {
	for _, u := range ur.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, store.ErrRecordNotFound
}
