package teststore

import (
	"github.com/Aza-9798/costs-rest-api/internal/app/model"
	"github.com/Aza-9798/costs-rest-api/internal/app/store"
)

type Store struct {
	userRepository        *UserRepository
	accountRepository     *AccountRepository
	transactionRepository *TransactionRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) User() store.UserRepo {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
			users: make(map[int]*model.User),
		}
	}
	return s.userRepository
}

func (s *Store) Account() store.AccountRepo {
	if s.accountRepository == nil {
		s.accountRepository = &AccountRepository{
			store:    s,
			accounts: make(map[int]*model.Account),
		}
	}
	return s.accountRepository
}

func (s *Store) Transaction() store.TransactionRepo {
	if s.transactionRepository == nil {
		s.transactionRepository = &TransactionRepository{
			store:        s,
			transactions: make(map[int]*model.TransactionDB),
		}
	}
	return s.transactionRepository
}
