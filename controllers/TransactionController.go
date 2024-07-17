package controllers

import (
	constants "TxnManagement/contants"
	"TxnManagement/controllers/exceptions"
	controllerModels "TxnManagement/controllers/models"
	"TxnManagement/controllers/utils"
	repoModels "TxnManagement/repositories/models"
	"TxnManagement/services"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
)

type TransactionController struct {
	customerService    *services.CustomerService
	transactionService *services.TransactionService
	transactionRetries int
	apiSecret          string
	logger             *constants.LoggingUtils
}

func NewTransactionController(customerRepository repoModels.CustomerRepository,
	transactionRepository repoModels.TransactionRepository,
	transactionRetries int,
	apiSecret string,
	logLevel string,
) *TransactionController {
	return &TransactionController{
		customerService:    services.NewCustomerService(customerRepository, logLevel),
		transactionService: services.NewTransactionService(transactionRepository, customerRepository, logLevel),
		transactionRetries: transactionRetries,
		apiSecret:          apiSecret,
		logger:             constants.NewLoggingUtils("TRANSACTION_CONTROLLER", logLevel),
	}
}

// AddTransaction godoc
// @Summary Adds a new transaction in DB
// @Description Registration of a new transaction for 1 product and 1 customer
// @Tags transaction
// @Accept json
// @Produce string
// @Success 200 {string}
// @Router /transaction [post]
func (t *TransactionController) AddTransaction(c *gin.Context) {
	email, err := utils.ExtractUsername(c, t.apiSecret)
	if err != nil {
		t.logger.Error("Error in admin authentication", err)
		c.JSON(401, err.Error())
		return
	}
	t.logger.Printf("Logged in admin: %s. Adding transaction", email)

	body, _ := io.ReadAll(c.Request.Body)
	var request controllerModels.TransactionRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		t.logger.Error("request not parseable", err)
		c.JSON(400, "Request not parseable")
		return
	}

	err = t.transactionService.AddTransaction(request)
	var noSuchCombinationError *exceptions.NoSuchCombinationError
	if errors.As(err, &noSuchCombinationError) {
		t.logger.Error("invalid request", err)
		c.JSON(400, "invalid request")
		return
	}
	if err != nil {
		t.logger.Error("error in adding transaction", err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in adding transaction. contact administrators for more info", HttpCode: 500})
		return
	}
	c.JSON(201, "Created transaction successfully")
}

// GetTransactions godoc
// @Summary Gets transactions from DB
// @Description Returns a list of transactions based on search criteria
// @Tags transaction
// @Accept json
// @Produce json
// @Success 200 {object} []controllerModels.TransactionResponse
// @Router /transaction [get]
func (t *TransactionController) GetTransactions(c *gin.Context) {
	email, err := utils.ExtractUsername(c, t.apiSecret)
	if err != nil {
		t.logger.Error("Error in admin authentication", err)
		c.JSON(401, err.Error())
		return
	}
	t.logger.Printf("Logged in admin: %s. Fetching transactions", email)
	values := c.Request.URL.Query()
	parsedParams, err := utils.ParseStrings(values, "customer_id", "village", "start_date", "end_date", "page", "page_size")
	if err != nil {
		t.logger.Error("Error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}

	t.logger.Printf("Fetching transactions for filter: "+
		"customer_id: %s, village:%s, start_date: %s, end_date: %s, page: %s, page_size: %s",
		utils.ToString(parsedParams["customer_id"]), utils.ToString(parsedParams["village"]),
		utils.ToString(parsedParams["start_date"]), utils.ToString(parsedParams["end_date"]),
		utils.ToString(parsedParams["page"]), utils.ToString(parsedParams["page_size"]))

	response, err := t.transactionService.FetchTransactions(parsedParams["customer_id"], parsedParams["village"],
		parsedParams["start_date"], parsedParams["end_date"], parsedParams["page"], parsedParams["page_size"])
	if err != nil {
		msg := fmt.Sprintf("error in gathering transactions for filter: "+
			"customer_id: %s, village:%s, start_date: %s, end_date: %s, page: %s, page_size: %s",
			utils.ToString(parsedParams["customer_id"]), utils.ToString(parsedParams["village"]),
			utils.ToString(parsedParams["start_date"]), utils.ToString(parsedParams["end_date"]),
			utils.ToString(parsedParams["page"]), utils.ToString(parsedParams["page_size"]))
		t.logger.Error(msg, err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting transactions. contact administrators for more info", HttpCode: 500})
		return
	}
	if response == nil {
		msg := fmt.Sprintf("transactions not found for filter: "+
			"customer_id: %s, village:%s, start_date: %s, end_date: %s, page: %s, page_size: %s",
			utils.ToString(parsedParams["customer_id"]), utils.ToString(parsedParams["village"]),
			utils.ToString(parsedParams["start_date"]), utils.ToString(parsedParams["end_date"]),
			utils.ToString(parsedParams["page"]), utils.ToString(parsedParams["page_size"]))
		t.logger.Info(msg)
		c.JSON(204, nil)
		return
	}
	c.JSON(200, response)
}

// GetCustomers godoc
// @Summary Gets customers from DB
// @Description Returns a list of customers based on search criteria
// @Tags customer
// @Accept json
// @Produce json
// @Success 200 {object} []models.CustomerData
// @Router /customer [get]
func (t *TransactionController) GetCustomers(c *gin.Context) {
	email, err := utils.ExtractUsername(c, t.apiSecret)
	if err != nil {
		t.logger.Error("error in admin authentication", err)
		c.JSON(401, err.Error())
		return
	}
	t.logger.Printf("Logged in admin: %s. Fetching customers", email)
	values := c.Request.URL.Query()
	parsedParams, err := utils.ParseStrings(values, "keyword", "page", "page_size")
	if err != nil {
		t.logger.Error("error in parsing params", err)
		c.JSON(400, err.Error())
		return
	}
	if parsedParams["keyword"] == nil {
		t.logger.Error("error in parsing keyword", err)
		c.JSON(400, errors.New("error in parsing keyword"))
		return
	}

	t.logger.Printf("Fetching customers for filter: keyword: %s, page: %s, page_size: %s",
		utils.ToString(parsedParams["keyword"]), utils.ToString(parsedParams["page"]), utils.ToString(parsedParams["page_size"]))

	response, err := t.customerService.FetchCustomers(*parsedParams["keyword"], parsedParams["page"], parsedParams["page_size"])
	if err != nil {
		msg := fmt.Sprintf("error in fetching customers for filter: keyword: %s, page: %s, page_size: %s",
			utils.ToString(parsedParams["keyword"]), utils.ToString(parsedParams["page"]), utils.ToString(parsedParams["page_size"]))
		t.logger.Error(msg, err)
		c.JSON(500, controllerModels.ErrorResponse{Message: "error in getting customers. contact administrators for more info", HttpCode: 500})
		return
	}
	if response == nil {
		msg := fmt.Sprintf("customers not found for filter: keyword: %s, page: %s, page_size: %s",
			utils.ToString(parsedParams["keyword"]), utils.ToString(parsedParams["page"]), utils.ToString(parsedParams["page_size"]))
		t.logger.Info(msg)
		c.JSON(204, nil)
		return
	}
	c.JSON(200, response)
}
