package sqlstore_test

import (
	"testing"

	"github.com/Aza-9798/costs-rest-api/internal/app/model"
	"github.com/Aza-9798/costs-rest-api/internal/app/store/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestAccountRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "accounts")
	store := sqlstore.New(db)

	u := model.TestUser(t)
	err := store.User().Create(u)
	assert.NoError(t, err)

	a := model.TestAccount(t, u)
	assert.NoError(t, store.Account().Create(a))
	assert.NotNil(t, a)
}

func TestAccountRepository_Find(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users", "accounts")
	store := sqlstore.New(db)

	u := model.TestUser(t)
	err := store.User().Create(u)
	assert.NoError(t, err)

	a := model.TestAccount(t, u)
	store.Account().Create(a)

	a1, err := store.Account().Find(a.ID)
	assert.NoError(t, err)
	assert.NotNil(t, a1)
}
