package models

import (
	"TxnManagement/repositories/models"
	"strings"
	"time"
)

type ErrorResponse struct {
	Message  string `json:"message"`
	HttpCode int    `json:"http_code"`
}

type TransactionResponse struct {
	CustomerName    string    `json:"customer_name"`
	CustomerFather  string    `json:"customer_father"`
	CustomerMobiles string    `json:"customer_mobiles"`
	CustomerAddress string    `json:"customer_address"`
	ProductName     string    `json:"product_name"`
	ProductType     string    `json:"product_type"`
	UnitPrice       int       `json:"unit_price"`
	TotalPrice      int       `json:"total_price"`
	AmountPaid      int       `json:"amount_paid"`
	AmountRemaining int       `json:"amount_remaining"`
	Date            time.Time `json:"date"`
}

func GetTransactionResponses(transactionData []models.TransactionData, customerData []models.CustomerData) []TransactionResponse {
	var customerIdMap map[string]models.CustomerData
	for _, customer := range customerData {
		customerIdMap[customer.Id] = customer
	}
	var response []TransactionResponse

	for _, transaction := range transactionData {
		allAddressSubs := customerIdMap[transaction.Customer.Id].Address.Tags
		allAddressSubs = append(allAddressSubs, customerIdMap[transaction.Customer.Id].Address.Village)
		response = append(response, TransactionResponse{
			CustomerName:    customerIdMap[transaction.Customer.Id].Name,
			CustomerFather:  customerIdMap[transaction.Customer.Id].Father,
			CustomerMobiles: strings.Join(customerIdMap[transaction.Customer.Id].Mobiles, ","),
			CustomerAddress: strings.Join(allAddressSubs, ","),
			ProductName:     strings.Join(transaction.Product.Tags, " "),
			ProductType:     transaction.Product.Type,
			UnitPrice:       transaction.Product.UnitPrice,
			TotalPrice:      transaction.Product.TotalPrice,
			AmountPaid:      transaction.AmountPaid,
			AmountRemaining: transaction.Product.TotalPrice - transaction.AmountPaid,
			Date:            transaction.Date,
		})
	}
	return response
}
