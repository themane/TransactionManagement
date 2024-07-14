package services

import (
	constants "TxnManagement/contants"
	"TxnManagement/repositories/models"
	"strconv"
)

type CustomerService struct {
	customerRepository models.CustomerRepository
	logger             *constants.LoggingUtils
}

func NewCustomerService(
	customerRepository models.CustomerRepository,
	logLevel string,
) *CustomerService {
	return &CustomerService{
		customerRepository: customerRepository,
		logger:             constants.NewLoggingUtils("CUSTOMER_SERVICE", logLevel),
	}
}

func (t *CustomerService) FetchCustomers(
	keyword string, pageString *string, pageSizeString *string) ([]models.CustomerData, error) {

	var err error
	page := 1
	pageSize := 10
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
	return t.customerRepository.FindByKeyword(keyword, page, pageSize)
}
