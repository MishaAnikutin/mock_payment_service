package repo

import (
	"context"
	"database/sql"
	"errors"
	"log"

	d "example.com/m/src/domain"
)

type PaymentRepo struct {
	db *sql.DB
}

func NewTransferRepo(db *sql.DB) d.PaymentRepoI {
	return &PaymentRepo{db: db}
}

func (r *PaymentRepo) Create(
	ctx context.Context,
	tx *sql.Tx,
	senderID string,
	receiverID string,
	amount int64,
) (*d.Payment, error) {
	payment := &d.Payment{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Amount:     amount,
		Status:     d.TransferPending,
	}

	const query = `
	INSERT INTO transfers (
		sender_id,
		receiver_id, 
		amount, 
		status
	) VALUES (?, ?, ?, ?)
	RETURNING id`

	err := tx.QueryRowContext(ctx, query,
		senderID,
		receiverID,
		amount,
		payment.Status,
	).Scan(&payment.ID)

	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *PaymentRepo) FindByID(
	ctx context.Context,
	transactionID int,
) (*d.Payment, error) {
	const query = `
		SELECT 
			id, 
			sender_id, 
			receiver_id, 
			amount, 
			status
		FROM 
			transfers 
		WHERE 
			id = ?
	`

	payment := &d.Payment{}
	err := r.db.QueryRowContext(ctx, query, transactionID).Scan(
		&payment.ID,
		&payment.SenderID,
		&payment.ReceiverID,
		&payment.Amount,
		&payment.Status,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, d.ErrPaymentNotFound
	case err != nil:
		return nil, d.ErrPaymentNotFound
	}

	return payment, nil
}

func (r *PaymentRepo) IsExists(ctx context.Context, paymentID int) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM transfers WHERE id = ?)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, paymentID).Scan(&exists)

	switch {
	case err != nil:
		return false, err
	case !exists:
		return false, d.ErrAccountNotFound
	}

	return true, nil
}

func (r *PaymentRepo) UpdateStatus(
	ctx context.Context,
	tx *sql.Tx,
	paymentID int,
	newStatus d.TransferStatus,
	expectedCurrentStatus d.TransferStatus,
) (*d.Payment, error) {
	const query = `
		UPDATE transfers 
		SET status = ? 
		WHERE id = ? AND status = ?
		RETURNING *;
	`

	log.Println("[cancel] существует ли:")

	payment := &d.Payment{}

	err := tx.QueryRowContext(ctx, query,
		newStatus,
		paymentID,
		expectedCurrentStatus,
	).Scan(
		&payment.ID,
		&payment.SenderID,
		&payment.ReceiverID,
		&payment.Amount,
		&payment.Status,
	)

	if err != nil {
		return payment, d.ErrPaymentIncorrectStatus
	}

	return payment, nil
}
