package sqlstore

import (
	"database/sql"
	"time"

	"github.com/Aza-9798/costs-rest-api/internal/app/model"
	"github.com/Aza-9798/costs-rest-api/internal/app/store"
)

type TransactionRepository struct {
	store *Store
}

func (r *TransactionRepository) Create(t *model.TransactionDB) error {
	if err := t.Validate(); err != nil {
		return err
	}
	//Create DB transaction for this operation
	tx, err := r.store.db.Begin()
	if err != nil {
		return err
	}

	if t.Type == model.StandardTransaction || t.Type == model.ExpenseTransaction {
		//Check account balance
		var balance float64
		if err := r.store.db.QueryRow(
			"select balance from accounts where id = $1",
			t.Source.ID,
		).Scan(&balance); err != nil {
			return err
		}
		if balance < t.Amount {
			return store.ErrInsufficientFunds
		}

		if _, err := tx.Exec("update accounts "+
			"set balance = balance - $1 "+
			"where id = $2",
			t.Amount,
			t.Source.ID); err != nil {
			tx.Rollback()
			return err
		}
	}

	if t.Type == model.IncomeTransaction || t.Type == model.StandardTransaction {
		if _, err := tx.Exec("update accounts "+
			"set balance = balance + $1 "+
			"where id = $2",
			t.Amount,
			t.Destination.ID); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.QueryRow("insert into transactions(transaction_date, source, destination, amount, description, type) "+
		"values($1, $2, $3, $4, $5, $6)"+
		"returning id, creation_date",
		t.TransactionDate,
		t.Source.ID,
		t.Destination.ID,
		t.Amount,
		t.Description,
		t.Type,
	).Scan(
		&t.ID,
		&t.CreationDate,
	); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *TransactionRepository) Delete(t *model.TransactionJSON) error {
	tx, err := r.store.db.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec("delete from transactions where id = $1", t.ID)
	if err != nil {
		return err
	}
	if count, err := res.RowsAffected(); err != nil {
		return err
	} else if count == 0 {
		return store.ErrRecordNotFound
	}
	if _, err := tx.Exec(
		"update accounts set balance = balance + $1 where id = $2",
		t.Amount,
		t.Source,
	); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(
		"update accounts set balance = balance - $1 where id = $2",
		t.Amount,
		t.Destination,
	); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *TransactionRepository) Save(t *model.TransactionDB) error {
	tInDB, err := r.Find(t.ID)
	if err != nil {
		return err
	}
	diff := t.Amount - tInDB.Amount
	if t.Type == model.ExpenseTransaction || t.Type == model.StandardTransaction {
		//изменить сумму у источника
		if t.Source.ID != tInDB.Source {
			t.Source.Balance -= t.Amount
			oldAcc, err := r.store.Account().Find(tInDB.Source)
			if err != nil {
				return err
			}
			oldAcc.Balance += tInDB.Amount
			r.store.Account().Save(oldAcc)
		} else if diff != 0 {
			t.Source.Balance -= diff
		}
		r.store.Account().Save(t.Source)
	}
	if t.Type == model.IncomeTransaction || t.Type == model.StandardTransaction {
		//изменить сумму у получателя
		if t.Destination.ID != tInDB.Destination {
			t.Destination.Balance += t.Amount
			oldAcc, err := r.store.Account().Find(tInDB.Destination)
			if err != nil {
				return err
			}
			oldAcc.Balance -= tInDB.Amount
			r.store.Account().Save(oldAcc)
		} else if diff != 0 {
			t.Destination.Balance += diff
		}
		r.store.Account().Save(t.Destination)
	}
	if _, err := r.store.db.Exec(
		"update transactions"+
			" set transaction_date = $1, source = $2, destination = $3, amount = $4, description = $5"+
			" where id = $6",
		t.TransactionDate,
		t.Source.ID,
		t.Destination.ID,
		t.Amount,
		t.Description,
		t.ID,
	); err != nil {
		return err
	}
	return nil
}

func (r *TransactionRepository) Find(id int) (*model.TransactionJSON, error) {
	t := &model.TransactionJSON{}
	if err := r.store.db.QueryRow("select id, creation_date, transaction_date, source, destination, amount, description"+
		" from transactions "+
		" where id = $1",
		id,
	).Scan(&t.ID,
		&t.CreationDate,
		&t.TransactionDate,
		&t.Source,
		&t.Destination,
		&t.Amount,
		&t.Description,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return t, nil
}

func (r *TransactionRepository) GetAllByAccount(accountID int) ([]*model.TransactionJSON, error) {
	rows, err := r.store.db.Query(
		"select id, creation_date, transaction_date, source, destination, amount, type, description"+
			" from transactions"+
			" where source = $1 or destination = $1", accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]*model.TransactionJSON, 0)
	for rows.Next() {
		t := &model.TransactionJSON{}
		err := rows.Scan(
			&t.ID,
			&t.CreationDate,
			&t.TransactionDate,
			&t.Source,
			&t.Destination,
			&t.Amount,
			&t.Type,
			&t.Description,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *TransactionRepository) GetAllByAccountAndPeriod(accountID int, DateStart, DateEnd time.Time) ([]*model.TransactionJSON, error) {
	rows, err := r.store.db.Query(
		"select id, creation_date, transaction_date, source, destination, amount, type, description"+
			" from transactions"+
			" where source = $1 or destination = $1"+
			" and transaction_date >= $2 and transaction_date <= $3",
		accountID,
		DateStart,
		DateEnd,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]*model.TransactionJSON, 0)
	for rows.Next() {
		t := &model.TransactionJSON{}
		err := rows.Scan(
			&t.ID,
			&t.CreationDate,
			&t.TransactionDate,
			&t.Source,
			&t.Destination,
			&t.Amount,
			&t.Type,
			&t.Description,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *TransactionRepository) GetAllByUser(userID int) ([]*model.TransactionJSON, error) {
	rows, err := r.store.db.Query(
		"select id, creation_date, transaction_date, source, destination, amount, type, description"+
			" from transactions "+
			" where source in (select id from accounts where user_id = $1)"+
			" or destination in (select id from accounts where user_id = $1)",
		userID)
	if err != nil {
		return nil, err
	}
	res := make([]*model.TransactionJSON, 0)
	for rows.Next() {
		t := &model.TransactionJSON{}
		err := rows.Scan(
			&t.ID,
			&t.CreationDate,
			&t.TransactionDate,
			&t.Source,
			&t.Destination,
			&t.Amount,
			&t.Type,
			&t.Description,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *TransactionRepository) GetSummary(userID int, DateStart, DateEnd time.Time) (*model.Summary, error) {
	res := &model.Summary{
		DateStart: DateStart,
		DateEnd:   DateEnd,
	}
	if err := r.store.db.QueryRow("select coalesce(sum(amount), 0) income"+
		" from transactions t, accounts a"+
		" where t.destination = a.id and a.user_id = $1"+
		" and type=$2"+
		" and t.transaction_date >= $3 and t.transaction_date <= $4",
		userID,
		model.IncomeTransaction,
		DateStart,
		DateEnd,
	).Scan(&res.Income); err != nil {
		return nil, err
	}
	if err := r.store.db.QueryRow("select coalesce(sum(amount), 0) expense"+
		" from transactions t, accounts a"+
		" where t.source = a.id and a.user_id = $1"+
		" and type=$2"+
		" and t.transaction_date >= $3 and t.transaction_date <= $4",
		userID,
		model.ExpenseTransaction,
		DateStart,
		DateEnd,
	).Scan(&res.Expense); err != nil {
		return nil, err
	}
	return res, nil
}
