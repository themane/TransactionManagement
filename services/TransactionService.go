package services

import (
	constants "TxnManagement/contants"
	"TxnManagement/controllers/exceptions"
	controllerModels "TxnManagement/controllers/models"
	"TxnManagement/repositories/models"
	"TxnManagement/utils"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type TransactionService struct {
	transactionRepository models.TransactionRepository
	customerRepository    models.CustomerRepository
	logger                *constants.LoggingUtils
}

func NewTransactionService(
	transactionRepository models.TransactionRepository,
	customerRepository models.CustomerRepository,
	logLevel string,
) *TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
		customerRepository:    customerRepository,
		logger:                constants.NewLoggingUtils("TRANSACTION_SERVICE", logLevel),
	}
}

func (t *TransactionService) AddTransaction(transactionRequest controllerModels.TransactionRequest) error {
	if !transactionRequest.Customer.New && (len(transactionRequest.Customer.Id) == 0 || len(transactionRequest.Customer.Address.Village) == 0) {
		return &exceptions.NoSuchCombinationError{Message: "invalid customer data"}
	}
	if transactionRequest.Customer.New {
		newUUID, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		transactionRequest.Customer.Id = newUUID.String()
		err = t.customerRepository.AddUser(transactionRequest.Customer)
		if err != nil {
			return err
		}
	}
	transactionData, err := transactionRequest.GetTransactionData()
	if err != nil {
		return err
	}
	err = t.transactionRepository.AddTransaction(*transactionData)
	if err != nil {
		return err
	}
	return nil
}

func (t *TransactionService) FetchTransactions(
	customerId *string, village *string, startDateString *string, endDateString *string, pageString *string, pageSizeString *string) ([]controllerModels.TransactionResponse, error) {

	var err error
	page := 1
	pageSize := 1000
	startDate := time.Now().AddDate(-1, 0, 0)
	endDate := time.Now()
	if pageString != nil && pageSizeString != nil {
		page, err = strconv.Atoi(*pageString)
		if err != nil {
			return nil, err
		}
		pageSize, err = strconv.Atoi(*pageSizeString)
		if err != nil {
			return nil, err
		}
	}
	if startDateString != nil && endDateString != nil {
		startDate, err = time.Parse(time.DateTime, *startDateString)
		if err != nil {
			return nil, err
		}
		endDate, err = time.Parse(time.DateTime, *endDateString)
		if err != nil {
			return nil, err
		}
	}

	if customerId != nil {
		transactionData, err := t.transactionRepository.FindByCustomerId(*customerId, startDate, endDate, page, pageSize)
		if err != nil {
			return nil, err
		}
		return t.getTransactionsResponse(transactionData)
	}
	if village != nil {
		transactionData, err := t.transactionRepository.FindByVillage(*village, startDate, endDate, page, pageSize)
		if err != nil {
			return nil, err
		}
		return t.getTransactionsResponse(transactionData)
	}
	transactionData, err := t.transactionRepository.FindByDate(startDate, endDate, page, pageSize)
	if err != nil {
		return nil, err
	}
	return t.getTransactionsResponse(transactionData)
}

func (t *TransactionService) getTransactionsResponse(transactionData []models.TransactionData) ([]controllerModels.TransactionResponse, error) {
	customerIdSet := utils.NewSet()
	for _, transaction := range transactionData {
		customerIdSet.Push(transaction.Customer.Id)
	}
	customerData, err := t.customerRepository.FindByIds(customerIdSet.Array())
	if err != nil {
		return nil, err
	}
	return controllerModels.GetTransactionResponses(transactionData, customerData), nil
}
