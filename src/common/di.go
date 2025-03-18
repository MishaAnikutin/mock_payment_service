package common

import (
	"database/sql"

	"example.com/m/src/application"
	"example.com/m/src/common/uow"
	"example.com/m/src/domain"
	"example.com/m/src/infrastructure"
	"example.com/m/src/infrastructure/repo"
	"example.com/m/src/presentation"
)

type Dependencies struct {
	DB               *sql.DB
	UoWFactory       uow.UnitOfWorkFactory
	AccountRepo      domain.AccountRepoI
	TransferRepo     domain.PaymentRepoI
	TransferUC       *application.TransferUC
	TransferHandlers *presentation.TransferHandlers
}

func Inject() (*Dependencies, error) {
	// session
	db, err := infrastructure.GetSession()
	if err != nil {
		return nil, err
	}

	// common
	uow_factory := uow.NewSqlUnitOfWorkFactory(db)
	// infra
	accountRepo := repo.NewAccountRepo(db)
	transferRepo := repo.NewTransferRepo(db)
	// application
	transferUC := application.NewTransferUC(accountRepo, transferRepo, uow_factory)
	// presentation
	transferHandlers := presentation.NewTransferHandlers(transferUC)

	return &Dependencies{
		DB:               db,
		UoWFactory:       uow_factory,
		AccountRepo:      accountRepo,
		TransferRepo:     transferRepo,
		TransferUC:       transferUC,
		TransferHandlers: transferHandlers,
	}, nil
}
