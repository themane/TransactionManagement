package models

import (
	controllerModels "TxnManagement/controllers/models"
	"github.com/google/uuid"
	"time"
)

type TransactionData struct {
	Id         string    `json:"-" bson:"_id"`
	Customer   Customer  `json:"customer" bson:"customer"`
	Product    Product   `json:"product" bson:"product"`
	AmountPaid int       `json:"amount_paid" bson:"amount_paid"`
	Date       time.Time `json:"date" bson:"date"`
}

type Customer struct {
	Id      string `json:"id" bson:"id"`
	Village string `json:"address" bson:"address"`
}

type Product struct {
	Type       string   `json:"type" bson:"type"`
	Weight     float32  `json:"weight" bson:"weight"`
	Tags       []string `json:"tags" bson:"tags"`
	UnitPrice  int      `json:"unit_price" bson:"unit_price"`
	TotalPrice int      `json:"total_price" bson:"total_price"`
}

func NewTransactionData(transactionRequest controllerModels.TransactionRequest) (*TransactionData, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &TransactionData{
		Id: newUUID.String(),
		Customer: Customer{
			Id:      transactionRequest.Customer.Id,
			Village: transactionRequest.Customer.Address.Village,
		},
		Product:    transactionRequest.Product,
		AmountPaid: transactionRequest.AmountPaid,
		Date:       time.Now(),
	}, nil
}

type TransactionRepository interface {
	FindById(id string) (*TransactionData, error)
	FindByCustomerId(customerId string, startDate time.Time, endDate time.Time, page int, pageSize int) ([]TransactionData, error)
	FindByVillage(village string, startDate time.Time, endDate time.Time, page int, pageSize int) ([]TransactionData, error)
	FindByDate(startDate time.Time, endDate time.Time, page int, pageSize int) ([]TransactionData, error)

	AddTransaction(transactionData TransactionData) error
}
