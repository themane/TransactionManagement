package models

import (
	"TxnManagement/models"
	"time"
)

type TransactionData struct {
	Id         string         `json:"-" bson:"_id"`
	Customer   Customer       `json:"customer" bson:"customer"`
	Product    models.Product `json:"product" bson:"product"`
	AmountPaid int            `json:"amount_paid" bson:"amount_paid"`
	Date       time.Time      `json:"date" bson:"date"`
}

type Customer struct {
	Id      string `json:"id" bson:"id"`
	Village string `json:"address" bson:"address"`
}

type TransactionRepository interface {
	FindById(id string) (*TransactionData, error)
	FindByCustomerId(customerId string, startDate time.Time, endDate time.Time, page int, pageSize int) ([]TransactionData, error)
	FindByVillage(village string, startDate time.Time, endDate time.Time, page int, pageSize int) ([]TransactionData, error)
	FindByDate(startDate time.Time, endDate time.Time, page int, pageSize int) ([]TransactionData, error)

	AddTransaction(transactionData TransactionData) error
}
