package model

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

func requieredIf(cond bool) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}

		return nil
	}
}

func validateTransactionType(transactionType, sourceAccountType, destinationAccountType string) validation.RuleFunc {
	return func(value interface{}) error {
		validationError := errors.New("transaction and account type mismatch")
		isStandardAccount := func(val string) bool {
			var standardAccountTypes []string = []string{CurrentAccount, DebtAccount, SavingAccount}
			for _, v := range standardAccountTypes {
				if v == val {
					return true
				}
			}
			return false
		}
		switch transactionType {
		case StandardTransaction:
			if !isStandardAccount(sourceAccountType) || !isStandardAccount(destinationAccountType) {
				return validationError
			}
		case IncomeTransaction:
			if !isStandardAccount(destinationAccountType) || sourceAccountType != IncomeSourceAccount {
				return validationError
			}
		case ExpenseTransaction:
			if !isStandardAccount(sourceAccountType) || destinationAccountType != ExpenseCatogoryAccount {
				return validationError
			}
		}
		return nil
	}
}
