package application

import (
	"context"
	"log"
	"time"

	"database/sql"

	d "example.com/m/src/domain"

	"example.com/m/src/common/uow"
)

type TransferUC struct {
	accRepo    d.AccountRepoI
	trRepo     d.PaymentRepoI
	uowFactory uow.UnitOfWorkFactory
}

func NewTransferUC(
	accRepo d.AccountRepoI,
	trRepo d.PaymentRepoI,
	uowFactory uow.UnitOfWorkFactory,
) *TransferUC {

	return &TransferUC{accRepo: accRepo, trRepo: trRepo, uowFactory: uowFactory}
}

func (uc *TransferUC) TransferMoney(ctx context.Context, account *d.Account, receiver_number string, amount int64) (*d.Payment, error) {
	log.Println("Проверяем счет отправителя")

	is_correct := uc.accRepo.CheckAccount(ctx, account)
	if !is_correct {
		return nil, d.ErrAccountIncorrect
	}

	log.Println("Получаем счет получателя")
	receiver_account, _ := uc.accRepo.GetByNumber(ctx, receiver_number)
	log.Println(receiver_account)

	log.Println("Создаем платеж")
	payment := &d.Payment{
		SenderID:   account.Number,
		ReceiverID: receiver_account.Number,
		Amount:     amount,
	}

	log.Println("Проверяем его")
	if err := uc.validatePayment(ctx, payment); err != nil {
		return nil, err
	}

	log.Println("Выполняем обмен средств")
	return uow.Execute[d.Payment](uc.uowFactory, ctx, func(tx *sql.Tx) (*d.Payment, error) {
		log.Println("Создаем трансфер")
		transfer, err := uc.trRepo.Create(ctx, tx,
			payment.SenderID,
			payment.ReceiverID,
			payment.Amount,
		)

		if err != nil {
			log.Println("Не удалось создать трансфер:", err)
			return nil, err
		}

		log.Println("Списываем деньги с одного счета")
		err = uc.accRepo.Withdraw(ctx, tx, payment.SenderID, payment.Amount)

		if err != nil {
			log.Panicln("Не удалось списать деньги с одного счета:", err)
			return nil, err
		}

		log.Println("Добавляем их на другой")
		err = uc.accRepo.Deposit(ctx, tx, payment.ReceiverID, payment.Amount)

		if err != nil {
			log.Panicln("Не удалось добавить их на другой:", err)
			return nil, err
		}

		log.Println("Обновляем статус платежа на COMPLETED")
		return uc.trRepo.UpdateStatus(ctx, tx, transfer.ID, d.TransferCompleted, d.TransferPending)
	})
}

func (uc *TransferUC) GetStatus(ctx context.Context, transfer_id int) (*d.TransferStatus, error) {
	transfer, err := uc.trRepo.FindByID(ctx, transfer_id)

	switch {
	case err != nil:
		return nil, d.ErrPaymentNotFound
	default:
		return &transfer.Status, nil
	}
}

func (uc *TransferUC) Cancel(ctx context.Context, transfer_id int) (*d.Payment, error) {
	is_exists, err := uc.trRepo.IsExists(ctx, transfer_id)

	log.Println("[cancel] существует ли:", is_exists)

	switch {
	case !is_exists:
		return nil, d.ErrPaymentNotFound
	case err != nil:
		return nil, err
	default:
		return uow.Execute[d.Payment](
			uc.uowFactory, ctx, func(tx *sql.Tx) (*d.Payment, error) {
				return uc.trRepo.UpdateStatus(ctx, tx, transfer_id, d.TransferPending, d.TransferCancelled)
			},
		)
	}
}

func (uc *TransferUC) validatePayment(ctx context.Context, payment *d.Payment) error {
	log.Println("  Проверяем отправителя на существование и что данные верные")
	account, err := uc.accRepo.GetByNumber(ctx, payment.ReceiverID)
	if err != nil {
		return err
	}

	if is_verified := uc.accRepo.CheckAccount(ctx, account); !is_verified {
		return d.ErrAccountIncorrect
	}

	log.Println("  Проверяем получателя на существование")
	if err := uc.accRepo.IsExists(ctx, payment.ReceiverID); err != nil {
		return d.ErrAccountNotFound
	}

	log.Println("  Проверяем что хватает средств у отправителя")
	if err := uc.accRepo.IsEnoughFunds(ctx, payment.SenderID, payment.Amount); err != nil {
		return d.ErrNotEnougthFunds
	}

	log.Println("  Проверяем размер платежа, что он больше нуля")
	if payment.Amount <= 0 {
		return d.ErrInvalidAmount
	}

	log.Println("  Проверяем, что это не платеж себе самому")
	if payment.SenderID == payment.ReceiverID {
		return d.ErrSameAccount
	}

	log.Println("  Имитируем бурную деятельность ... тут можно всякие антифрод, KYC и тд")
	time.Sleep(3 * time.Second)
	return nil
}
