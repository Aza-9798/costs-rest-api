package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	StandardTransaction = "Standard"
	IncomeTransaction   = "Income"
	ExpenseTransaction  = "Expense"
)

type TransactionDB struct {
	ID              int       `json:"id"`
	CreationDate    time.Time `json:"-"`
	TransactionDate time.Time `json:"transaction_date"`
	Source          *Account  `json:"source"`
	Destination     *Account  `json:"destination"`
	Amount          float64   `json:"amount"`
	Type            string    `json:"type"`
	Description     string    `json:"description"`
}

type TransactionJSON struct {
	ID              int       `json:"id"`
	CreationDate    time.Time `json:"-"`
	TransactionDate time.Time `json:"transaction_date"`
	Source          int       `json:"source"`
	Destination     int       `json:"destination"`
	Amount          float64   `json:"amount"`
	Type            string    `json:"type"`
	Description     string    `json:"description"`
}

func (t *TransactionDB) Validate() error {
	return validation.ValidateStruct(
		t,
		validation.Field(&t.Source, validation.Required),
		validation.Field(&t.Destination, validation.Required),
		validation.Field(&t.Amount, validation.Required),
		validation.Field(&t.Type, validation.Required,
			validation.In(
				StandardTransaction,
				IncomeTransaction,
				ExpenseTransaction,
			),
			validation.By(
				validateTransactionType(
					t.Type,
					t.Source.Type,
					t.Destination.Type,
				),
			),
		),
	)
}

func (t *TransactionDB) ToJSON() *TransactionJSON {
	res := &TransactionJSON{
		ID:              t.ID,
		CreationDate:    t.CreationDate,
		TransactionDate: t.TransactionDate,
		Source:          t.Source.ID,
		Destination:     t.Destination.ID,
		Amount:          t.Amount,
		Type:            t.Type,
		Description:     t.Description,
	}
	return res
}
