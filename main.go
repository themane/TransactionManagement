package main

import (
	constants "TxnManagement/contants"
	"TxnManagement/controllers"
	"TxnManagement/repositories"
	secretManager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"sync"
)

var once = sync.Once{}
var mongoDB string
var mongoUrlSecretName string
var apiSecretName string
var transactionRetries int
var logLevel string

// @title Transaction Management Server
// @version 1.0.0
// @description This is the server for any Transaction Management system
// @termsOfService http://swagger.io/terms/

// @contact.name Devashish Gupta
// @contact.email devagpta@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @schemes https
func main() {
	r := gin.Default()

	once.Do(initialize)

	authController, transactionController := getHandlers()

	r.GET("/ping", controllers.Ping)
	r.POST("/admin/login", authController.Login)
	r.POST("/admin/register", authController.Register)
	r.HEAD("/admin/token", authController.RefreshToken)
	r.POST("/transaction", transactionController.AddTransaction)
	r.GET("/transaction", transactionController.GetTransactions)
	r.GET("/customer", transactionController.GetCustomers)

	err := r.Run()
	if err != nil {
		log.Println("Error in starting server")
		return
	}
}

func getHandlers() (*controllers.AuthController, *controllers.TransactionController) {
	log.Println("Initializing handlers")

	mongoURL := accessSecretVersion(mongoUrlSecretName)
	apiSecret := accessSecretVersion(apiSecretName)

	adminRepository := repositories.NewAdminRepository(mongoURL, mongoDB, logLevel)
	customerRepository := repositories.NewCustomerRepository(mongoURL, mongoDB, logLevel)
	transactionRepository := repositories.NewTransactionRepository(mongoURL, mongoDB, logLevel)

	authController := controllers.NewAuthController(adminRepository, apiSecret, logLevel)
	transactionController := controllers.NewTransactionController(customerRepository, transactionRepository,
		transactionRetries, apiSecret, logLevel)

	log.Println("Initialized all handlers")
	return authController, transactionController
}

func initialize() {
	mongoUrlSecretName = os.Getenv("MONGO_SECRET_NAME")
	mongoDB = os.Getenv("MONGO_DB")
	if mongoUrlSecretName == "" || mongoDB == "" {
		log.Fatal("Mongo not configured")
	}

	apiSecretName = os.Getenv("API_SECRET_NAME")
	if apiSecretName == "" {
		apiSecretName = "API_SECRET"
	}

	logLevel = os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = constants.Info
	}
}

func accessSecretVersion(secretName string) string {
	ctx := context.Background()
	client, err := secretManager.NewClient(ctx)
	if err != nil {
		log.Fatal("Error in initializing client for secret manager: ", err)
		return ""
	}
	defer func(client *secretManager.Client) {
		err := client.Close()
		if err != nil {
			log.Fatal("Error in closing client for secret manager: ", err)
		}
	}(client)
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		log.Fatal("Error in calling access API for retrieving secret data: ", err)
		return ""
	}
	return string(result.Payload.Data)
}
