package repo

import (
	"context"
	"database/sql"
	"errors"
	"log"

	d "example.com/m/src/domain"
)

type AccountRepo struct {
	db *sql.DB
}

func NewAccountRepo(db *sql.DB) d.AccountRepoI {
	return &AccountRepo{db: db}
}

func (r *AccountRepo) CheckAccount(ctx context.Context, account *d.Account) bool {
	log.Println("\tПолучаем данные счета по номеру")
	dbAccount, err := r.GetByNumber(ctx, account.Number)

	if err != nil {
		return false
	}
	is_ok := dbAccount.FullName == account.FullName &&
		dbAccount.CVV == account.CVV &&
		dbAccount.ExparationDate == account.ExparationDate

	log.Println("\tВ базе:", dbAccount, "Ввели:", account, "Итог:", is_ok)

	return is_ok
}

func (r *AccountRepo) GetByNumber(ctx context.Context, number string) (*d.Account, error) {
	const query = `
	SELECT number, exparation_date, full_name, cvv
	FROM accounts
	WHERE number = ?
	`

	log.Println("\tПолучаем данные")
	account := &d.Account{}

	err := r.db.QueryRowContext(ctx, query, number).Scan(
		&account.Number,
		&account.ExparationDate,
		&account.FullName,
		&account.CVV,
	)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, d.ErrAccountNotFound
	}

	return account, err
}

func (r *AccountRepo) IsExists(ctx context.Context, accountID string) error {
	const query = `SELECT EXISTS(SELECT 1 FROM accounts WHERE number = ?)`

	var exists bool

	err := r.db.QueryRowContext(ctx, query, accountID).Scan(&exists)

	switch {
	case err != nil:
		return err
	case !exists:
		return d.ErrAccountNotFound
	}

	return nil
}

func (r *AccountRepo) IsEnoughFunds(ctx context.Context, accountID string, amount int64) error {
	const query = `SELECT balance FROM accounts WHERE number = ?`

	var balance int64
	err := r.db.QueryRowContext(ctx, query, accountID).Scan(&balance)

	switch {
	case err == sql.ErrNoRows:
		return d.ErrAccountNotFound
	case err != nil:
		return err
	case balance < amount:
		return d.ErrNotEnoughFunds
	}

	return nil
}

func (r *AccountRepo) Deposit(ctx context.Context, tx *sql.Tx, accountID string, amount int64) error {
	const query = `
        UPDATE accounts
        SET balance = balance + ?
        WHERE number = ?
    `

	res, err := tx.ExecContext(ctx, query, amount, accountID)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 {
		return d.ErrAccountNotFound
	}

	return err
}

func (r *AccountRepo) Withdraw(ctx context.Context, tx *sql.Tx, accountID string, amount int64) error {
	const query = `
        UPDATE accounts 
        SET balance = balance - ?
        WHERE number = ?
    `

	res, err := tx.ExecContext(ctx, query, amount, accountID)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected == 0 {
		return d.ErrAccountNotFound
	}

	return err
}
