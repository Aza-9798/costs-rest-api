package model

import "testing"

func TestUser(t *testing.T) *User {
	t.Helper()
	return &User{
		Email:    "user@example.org",
		Password: "password",
	}
}

func TestAccount(t *testing.T, u *User) *Account {
	return &Account{
		User:        u.ID,
		Name:        "Test",
		Type:        CurrentAccount,
		Description: "Test",
		Balance:     100,
	}
}
