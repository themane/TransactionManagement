package models

import (
	"TxnManagement/models"
	repoModels "TxnManagement/repositories/models"
	"time"

	"github.com/google/uuid"
)

type TransactionRequest struct {
	Customer   repoModels.CustomerData `json:"customer"`
	Product    models.Product          `json:"product"`
	AmountPaid int                     `json:"amount_paid"`
}

func (t *TransactionRequest) GetTransactionData() (*repoModels.TransactionData, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &repoModels.TransactionData{
		Id: newUUID.String(),
		Customer: repoModels.Customer{
			Id:      t.Customer.Id,
			Village: t.Customer.Address.Village,
		},
		Product:    t.Product,
		AmountPaid: t.AmountPaid,
		Date:       time.Now(),
	}, nil
}
