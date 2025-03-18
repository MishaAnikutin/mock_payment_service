package uow

import (
	"context"

	"database/sql"
)

type UnitOfWork interface {
	Begin(ctx context.Context) (*sql.Tx, error)
	Commit() error
	Rollback() error
}

type UnitOfWorkFactory interface {
	New(ctx context.Context) (UnitOfWork, error)
}

func Execute[T any](
	factory UnitOfWorkFactory,
	ctx context.Context,
	fn func(tx *sql.Tx) (*T, error),
) (*T, error) {
	uow, err := factory.New(ctx)
	if err != nil {
		return nil, err
	}
	defer uow.Rollback()

	tx, err := uow.Begin(ctx)
	if err != nil {
		return nil, err
	}

	res, err := fn(tx)
	if err != nil {
		return nil, err
	}

	return res, uow.Commit()
}
