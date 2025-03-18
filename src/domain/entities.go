package domain

type Account struct {
	Number         string
	ExparationDate string
	FullName       string
	CVV            uint16
}

type TransferStatus string

const (
	TransferPending   TransferStatus = "PENDING"
	TransferCompleted TransferStatus = "COMPLETED"
	TransferFailed    TransferStatus = "FAILED"
	TransferCancelled TransferStatus = "CANCELLED"
)

type Payment struct {
	ID         int
	SenderID   string
	ReceiverID string
	Amount     int64
	Status     TransferStatus
}
