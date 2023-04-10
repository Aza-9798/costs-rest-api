package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	CurrentAccount         = "Current"
	SavingAccount          = "Saving"
	DebtAccount            = "Debt"
	IncomeSourceAccount    = "IncomeSource"
	ExpenseCatogoryAccount = "ExpenseCategory"
)

type Account struct {
	ID           int       `json:"id"`
	CreationDate time.Time `json:"-"`
	User         int       `json:"user"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Description  string    `json:"description"`
	Balance      float64   `json:"balance"`
}

func (a *Account) Validate() error {
	return validation.ValidateStruct(
		a,
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.Type,
			validation.Required,
			validation.In(CurrentAccount,
				SavingAccount,
				DebtAccount,
				IncomeSourceAccount,
				ExpenseCatogoryAccount,
			)),
	)
}
