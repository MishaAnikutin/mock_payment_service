package domain

import (
	"context"
	"database/sql"
)

type AccountRepoI interface {
	IsExists(ctx context.Context, accountID string) error

	CheckAccount(ctx context.Context, account *Account) bool

	GetByNumber(ctx context.Context, number string) (*Account, error)

	IsEnoughFunds(ctx context.Context, accoundID string, amount int64) error

	Deposit(ctx context.Context, tx *sql.Tx, accoundID string, amount int64) error

	Withdraw(ctx context.Context, tx *sql.Tx, accoundID string, amount int64) error
}

type PaymentRepoI interface {
	Create(
		ctx context.Context,
		tx *sql.Tx,
		SenderID string,
		ReceiverID string,
		Amount int64,
	) (*Payment, error)

	FindByID(
		ctx context.Context,
		transactionID int,
	) (*Payment, error)

	IsExists(
		ctx context.Context,
		transactionID int,
	) (bool, error)

	UpdateStatus(
		ctx context.Context,
		tx *sql.Tx,
		transactionID int,
		status TransferStatus,
		// Для консистентости данных добавляем ожидаемое состояние транзакции до обновления.
		// Позволит решить проблему Lost Update без акторной модели
		expectedStatus TransferStatus,
	) (*Payment, error)
}
