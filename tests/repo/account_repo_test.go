package repo_test

import (
	"context"
	"database/sql"
	"testing"

	d "example.com/m/src/domain"
	"example.com/m/src/infrastructure/repo"
	"example.com/m/tests/fixtures"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func getRepoAndMock(t *testing.T) (*sql.DB, d.AccountRepoI, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("не удалось создать мок БД: %s", err)
	}

	return db, repo.NewAccountRepo(db), mock
}

func TestAccountRepo_GetByNumber(t *testing.T) {
	db, r, mock := getRepoAndMock(t)
	defer db.Close()

	t.Run("Success", func(t *testing.T) {
		account := fixtures.ValidAccount

		rows := sqlmock.NewRows([]string{"number", "exparation_date", "full_name", "cvv"}).
			AddRow(account.Number, account.ExparationDate, account.FullName, account.CVV)

		mock.ExpectQuery("SELECT (.+) FROM accounts WHERE number = ?").
			WithArgs(account.Number).
			WillReturnRows(rows)

		account, err := r.GetByNumber(context.Background(), account.Number)

		assert.NoError(t, err)
		assert.Equal(t, &d.Account{
			Number:         account.Number,
			ExparationDate: account.ExparationDate,
			FullName:       account.FullName,
			CVV:            account.CVV,
		}, account)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM accounts WHERE number = ?").
			WithArgs(fixtures.InvalidAccountNumber).
			WillReturnError(sql.ErrNoRows)

		account, err := r.GetByNumber(context.Background(), fixtures.InvalidAccountNumber)

		assert.Nil(t, account)
		assert.ErrorIs(t, err, d.ErrAccountNotFound)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAccountRepo_CheckAccount(t *testing.T) {
	db, r, mock := getRepoAndMock(t)
	defer db.Close()

	t.Run("Success", func(t *testing.T) {
		account := fixtures.ValidAccount

		rows := sqlmock.NewRows([]string{"number", "exparation_date", "full_name", "cvv"}).
			AddRow(account.Number, account.ExparationDate, account.FullName, account.CVV)

		mock.ExpectQuery("SELECT (.+) FROM accounts WHERE number = ?").
			WithArgs(fixtures.ValidAccount.Number).
			WillReturnRows(rows)

		is_valid := r.CheckAccount(context.Background(), account)

		assert.True(t, is_valid)
	})

	t.Run("Not valid", func(t *testing.T) {
		account := fixtures.ValidAccount

		rows := sqlmock.NewRows([]string{"number", "exparation_date", "full_name", "cvv"}).
			AddRow(account.Number, account.ExparationDate, account.FullName, account.CVV)

		mock.ExpectQuery("SELECT (.+) FROM accounts WHERE number = ?").
			WithArgs(fixtures.ValidAccount.Number).
			WillReturnRows(rows)

		// мухухехе
		account.Number = fixtures.InvalidAccountNumber

		is_valid := r.CheckAccount(context.Background(), account)

		assert.False(t, is_valid)
	})

	t.Run("Not found", func(t *testing.T) {
		account := fixtures.ValidAccount
		account.Number = fixtures.InvalidAccountNumber

		rows := sqlmock.NewRows([]string{"number", "exparation_date", "full_name", "cvv"}).
			AddRow(account.Number, account.ExparationDate, account.FullName, account.CVV)

		mock.ExpectQuery("SELECT (.+) FROM accounts WHERE number = ?").
			WithArgs(fixtures.ValidAccount.Number).
			WillReturnRows(rows)

		is_valid := r.CheckAccount(context.Background(), account)

		assert.False(t, is_valid)
	})
}
