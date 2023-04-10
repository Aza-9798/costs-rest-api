package model

import (
	"errors"
	"time"
)

type Summary struct {
	DateStart time.Time `json:"date_start"`
	DateEnd   time.Time `json:"date_end"`
	Income    float64   `json:"income"`
	Expense   float64   `json:"expense"`
}

type AccountSummary interface {
	SetPeriod(time.Time, time.Time) error
	CalculateSummary([]*TransactionJSON)
}

type IncomeAccountSummary struct {
	AccountID int       `json:"account_id"`
	DateStart time.Time `json:"date_start"`
	DateEnd   time.Time `json:"date_end"`
	Income    float64   `json:"income"`
}

func validatePeriod(DateStart, DateEnd time.Time) error {
	if DateStart.After(DateEnd) || DateStart.Equal(DateEnd) {
		return errors.New("wrong or empty period")
	}
	return nil
}

func (s *IncomeAccountSummary) SetPeriod(DateStart, DateEnd time.Time) error {
	if err := validatePeriod(DateStart, DateEnd); err != nil {
		return err
	}
	s.DateStart, s.DateEnd = DateStart, DateEnd
	return nil
}

func (s *IncomeAccountSummary) CalculateSummary(ts []*TransactionJSON) {
	var income float64
	for _, t := range ts {
		if t.Source == s.AccountID {
			income += t.Amount
		}
	}
	s.Income = income
}

type ExpenseCategorySummary struct {
	AccountID int       `json:"account_id"`
	DateStart time.Time `json:"date_start"`
	DateEnd   time.Time `json:"date_end"`
	Expense   float64   `json:"expense"`
}

func (s *ExpenseCategorySummary) SetPeriod(DateStart, DateEnd time.Time) error {
	if err := validatePeriod(DateStart, DateEnd); err != nil {
		return err
	}
	s.DateStart, s.DateEnd = DateStart, DateEnd
	return nil
}

func (s *ExpenseCategorySummary) CalculateSummary(ts []*TransactionJSON) {
	var expense float64
	for _, t := range ts {
		if t.Destination == s.AccountID {
			expense += t.Amount
		}
	}
	s.Expense = expense
}

type StandardAccountSummary struct {
	AccountID int       `json:"account_id"`
	DateStart time.Time `json:"date_start"`
	DateEnd   time.Time `json:"date_end"`
	Income    float64   `json:"income"`
	Expense   float64   `json:"expense"`
}

func (s *StandardAccountSummary) SetPeriod(DateStart, DateEnd time.Time) error {
	if err := validatePeriod(DateStart, DateEnd); err != nil {
		return err
	}
	s.DateStart, s.DateEnd = DateStart, DateEnd
	return nil
}

func (s *StandardAccountSummary) CalculateSummary(ts []*TransactionJSON) {
	var income, expense float64
	for _, t := range ts {
		if t.Destination == s.AccountID {
			income += t.Amount
		}
		if t.Source == s.AccountID {
			expense += t.Amount
		}
	}
	s.Income, s.Expense = income, expense
}

func GetSummaryByAccount(a *Account) (AccountSummary, error) {
	if a.Type == CurrentAccount || a.Type == DebtAccount || a.Type == SavingAccount {
		return &StandardAccountSummary{AccountID: a.ID}, nil
	} else if a.Type == ExpenseCatogoryAccount {
		return &ExpenseCategorySummary{AccountID: a.ID}, nil
	} else if a.Type == IncomeSourceAccount {
		return &IncomeAccountSummary{AccountID: a.ID}, nil
	} else {
		return nil, errors.New("no AccountSummary for this account type")
	}
}
