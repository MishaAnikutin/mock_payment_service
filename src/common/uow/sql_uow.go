package uow

import (
	"context"
	"database/sql"
	"errors"
)

type SqlUnitOfWork struct {
	db *sql.DB
	tx *sql.Tx
}

func NewSqlUnitOfWorkFactory(db *sql.DB) UnitOfWorkFactory {
	return &sqlUowFactory{db: db}
}

type sqlUowFactory struct {
	db *sql.DB
}

func (f *sqlUowFactory) New(ctx context.Context) (UnitOfWork, error) {
	return &SqlUnitOfWork{db: f.db}, nil
}

func (u *SqlUnitOfWork) Begin(ctx context.Context) (*sql.Tx, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	u.tx = tx
	return tx, nil
}

func (u *SqlUnitOfWork) Commit() error {
	if u.tx == nil {
		return errors.New("no transaction to commit")
	}
	return u.tx.Commit()
}

func (u *SqlUnitOfWork) Rollback() error {
	if u.tx == nil {
		return nil
	}
	return u.tx.Rollback()
}
