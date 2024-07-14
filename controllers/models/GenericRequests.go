package models

import "TxnManagement/repositories/models"

type TransactionRequest struct {
	Customer   models.CustomerData `json:"customer"`
	Product    models.Product      `json:"product"`
	AmountPaid int                 `json:"amount_paid"`
}
