package dto

import "github.com/google/uuid"

type TransferRequest struct {
	FromAccountID uuid.UUID `json:"from_account_id" binding:"required"`
	ToAccountID   uuid.UUID `json:"to_account_id" binding:"required"`
	Amount        float64   `json:"amount" binding:"required,gt=0"`
	Description   string    `json:"description"`
}

type TransferResponse struct {
	TransactionID         uuid.UUID `json:"transaction_id"`
	TransactionReference  string    `json:"transaction_reference"`
	Status                string    `json:"status"`
	Amount                float64   `json:"amount"`
	CurrencyCode          string    `json:"currency_code"`
	FromAccountID         uuid.UUID `json:"from_account_id"`
	ToAccountID           uuid.UUID `json:"to_account_id"`
	FromAccountNewBalance float64   `json:"from_account_new_balance"`
	ToAccountNewBalance   float64   `json:"to_account_new_balance"`
}
