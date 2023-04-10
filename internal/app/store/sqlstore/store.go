package sqlstore

import (
	"database/sql"

	"github.com/Aza-9798/costs-rest-api/internal/app/store"
	_ "github.com/lib/pq"
)

type Store struct {
	db                    *sql.DB
	userRepository        *UserRepository
	accountRepository     *AccountRepository
	transactionRepository *TransactionRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepo {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			store: s,
		}
	}
	return s.userRepository
}

func (s *Store) Account() store.AccountRepo {
	if s.accountRepository == nil {
		s.accountRepository = &AccountRepository{
			store: s,
		}
	}
	return s.accountRepository
}

func (s *Store) Transaction() store.TransactionRepo {
	if s.transactionRepository == nil {
		s.transactionRepository = &TransactionRepository{
			store: s,
		}
	}
	return s.transactionRepository
}
