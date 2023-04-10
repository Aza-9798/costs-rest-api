package store

import (
	"time"

	"github.com/Aza-9798/costs-rest-api/internal/app/model"
)

type UserRepo interface {
	Create(user *model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}

type AccountRepo interface {
	Create(account *model.Account) error
	Delete(int) error
	Save(*model.Account) error
	Find(int) (*model.Account, error)
	GetAllByUser(int) ([]*model.Account, error)
	IsAccountBelongsUser(int, int) (bool, error)
}

type TransactionRepo interface {
	Create(transaction *model.TransactionDB) error
	Delete(transaction *model.TransactionJSON) error
	Save(*model.TransactionDB) error
	Find(int) (*model.TransactionJSON, error)
	GetAllByAccount(int) ([]*model.TransactionJSON, error)
	GetAllByAccountAndPeriod(int, time.Time, time.Time) ([]*model.TransactionJSON, error)
	GetAllByUser(int) ([]*model.TransactionJSON, error)
	GetSummary(int, time.Time, time.Time) (*model.Summary, error)
}
