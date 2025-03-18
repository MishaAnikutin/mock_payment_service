package presentation

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"example.com/m/src/application"
	"example.com/m/src/domain"

	"github.com/gin-gonic/gin"
)

type CreateTransferRequest struct {
	Number         string `json:"number" binding:"required"`
	FullName       string `json:"full_name" binding:"required"`
	YY             string `json:"yy" binding:"required"`
	MM             string `json:"mm" binding:"required"`
	Cvv            int    `json:"cvv" binding:"required"`
	ReceiverNumber string `json:"receiver_number" binding:"required"`
	Amount         int64  `json:"amount" binding:"required,gt=0"`
}

type CancelTransferRequest struct {
	TransferID int `json:"transfer_id" binding:"required"`
}

type TransferHandlers struct {
	transferUC *application.TransferUC
}

func NewTransferHandlers(uc *application.TransferUC) *TransferHandlers {
	return &TransferHandlers{transferUC: uc}
}

// POST /transfer
func (h *TransferHandlers) CreateTransfer(c *gin.Context) {
	var req CreateTransferRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account := &domain.Account{
		Number:         req.Number,
		FullName:       req.FullName,
		CVV:            uint16(req.Cvv),
		ExparationDate: fmt.Sprintf("%s-%s", req.MM, req.YY),
	}

	payment, err := h.transferUC.TransferMoney(c.Request.Context(), account, req.ReceiverNumber, req.Amount)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Transfer initiated successfully",
		"id":      payment.ID,
	})
}

// POST /transfer/cancel
func (h *TransferHandlers) CancelTransfer(c *gin.Context) {
	var req CancelTransferRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.transferUC.Cancel(c.Request.Context(), req.TransferID); err != nil {
		log.Println(err)
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transfer cancelled successfully",
	})
}

// GET /transfer/status/:id
func (h *TransferHandlers) GetStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transfer ID"})
		return
	}

	status, err := h.transferUC.GetStatus(c.Request.Context(), id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}

func handleError(c *gin.Context, err error) {
	switch err {
	case domain.ErrPaymentNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "Transfer not found"})
	case domain.ErrAccountIncorrect:
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
	case domain.ErrAccountNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
	case domain.ErrNotEnougthFunds:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
	case domain.ErrSameAccount:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot transfer to same account"})
	case domain.ErrInvalidAmount:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}
