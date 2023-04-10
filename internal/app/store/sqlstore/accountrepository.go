package sqlstore

import (
	"database/sql"

	"github.com/Aza-9798/costs-rest-api/internal/app/model"
	"github.com/Aza-9798/costs-rest-api/internal/app/store"
)

type AccountRepository struct {
	store *Store
}

func (r *AccountRepository) Create(a *model.Account) error {
	if err := a.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow("insert into accounts(name, user_id, account_type, description, balance) values($1, $2, $3, $4, $5) returning id, creation_date",
		a.Name,
		a.User,
		a.Type,
		a.Description,
		a.Balance,
	).Scan(&a.ID, &a.CreationDate)
}

func (r *AccountRepository) Save(a *model.Account) error {
	if err := a.Validate(); err != nil {
		return err
	}
	if _, err := r.Find(a.ID); err != nil {
		return store.ErrRecordNotFound
	}
	if _, err := r.store.db.Exec(
		"update accounts"+
			" set name = $1, description = $2, balance = $3"+
			" where id = $4",
		a.Name,
		a.Description,
		a.Balance,
		a.ID,
	); err != nil {
		return err
	}
	return nil
}

func (r *AccountRepository) Find(id int) (*model.Account, error) {
	a := &model.Account{}
	if err := r.store.db.QueryRow("select * from accounts where id = $1",
		id,
	).Scan(
		&a.ID,
		&a.CreationDate,
		&a.Type,
		&a.Balance,
		&a.User,
		&a.Name,
		&a.Description,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return a, nil
}

func (r *AccountRepository) Delete(id int) error {
	res, err := r.store.db.Exec("delete from accounts where id = $1", id)
	if err != nil {
		return err
	}
	if count, err := res.RowsAffected(); err != nil {
		return err
	} else if count == 0 {
		return store.ErrRecordNotFound
	}
	return nil
}

func (r *AccountRepository) GetAllByUser(userID int) ([]*model.Account, error) {
	rows, err := r.store.db.Query(
		"select id, creation_date, user_id, name, account_type, balance, description "+
			"from accounts "+
			"where user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]*model.Account, 0)
	for rows.Next() {
		a := &model.Account{}
		err := rows.Scan(
			&a.ID,
			&a.CreationDate,
			&a.User,
			&a.Name,
			&a.Type,
			&a.Balance,
			&a.Description,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *AccountRepository) IsAccountBelongsUser(accountID, userID int) (bool, error) {
	var existFlag int
	if err := r.store.db.QueryRow(
		"select case when exists(select * from accounts where id = $1 and user_id = $2) then 1 else 0 end",
		accountID,
		userID,
	).Scan(&existFlag); err != nil {
		return false, err
	}
	return existFlag == 1, nil
}
