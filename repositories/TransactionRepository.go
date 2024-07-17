package repositories

import (
	constants "TxnManagement/contants"
	"TxnManagement/repositories/models"
	"TxnManagement/repositories/utils"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionRepositoryImpl struct {
	mongoURL string
	mongoDB  string
	logger   *constants.LoggingUtils
}

func NewTransactionRepository(mongoURL string, mongoDB string, logLevel string) *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{
		mongoURL: mongoURL,
		mongoDB:  mongoDB,
		logger:   constants.NewLoggingUtils("TRANSACTION_REPOSITORY", logLevel),
	}
}

func (u *TransactionRepositoryImpl) FindById(id string) (*models.TransactionData, error) {
	client, ctx, cancel := u.getMongoClient()
	defer utils.Disconnect(client, ctx)
	filter := bson.M{"_id": id}
	return u.findSingleResult(ctx, u.getCollection(client), cancel, filter)
}

func (u *TransactionRepositoryImpl) FindByCustomerId(customerId string, startDate time.Time, endDate time.Time,
	page int, pageSize int) ([]models.TransactionData, error) {

	client, ctx, cancel := u.getMongoClient()
	defer utils.Disconnect(client, ctx)
	filter := bson.M{"customer.id": customerId, "date": bson.M{"$gte": startDate, "$lt": endDate}}
	return u.findResults(ctx, u.getCollection(client), cancel, filter, page, pageSize)
}

func (u *TransactionRepositoryImpl) FindByVillage(village string, startDate time.Time, endDate time.Time,
	page int, pageSize int) ([]models.TransactionData, error) {

	client, ctx, cancel := u.getMongoClient()
	defer utils.Disconnect(client, ctx)
	filter := bson.M{"customer.village": village, "date": bson.M{"$gte": startDate, "$lt": endDate}}
	return u.findResults(ctx, u.getCollection(client), cancel, filter, page, pageSize)
}

func (u *TransactionRepositoryImpl) FindByDate(startDate time.Time, endDate time.Time,
	page int, pageSize int) ([]models.TransactionData, error) {

	client, ctx, cancel := u.getMongoClient()
	defer utils.Disconnect(client, ctx)
	filter := bson.M{"date": bson.M{"$gte": startDate, "$lt": endDate}}
	return u.findResults(ctx, u.getCollection(client), cancel, filter, page, pageSize)
}

func (u *TransactionRepositoryImpl) AddTransaction(transactionData models.TransactionData) error {
	client, ctx, cancel := u.getMongoClient()
	defer utils.Disconnect(client, ctx)
	_, err := u.getCollection(client).InsertOne(ctx, transactionData)
	defer cancel()
	if err != nil {
		u.logger.Error("Error in saving transaction data to Mongo", err)
		return err
	}
	return nil
}

func (u *TransactionRepositoryImpl) getMongoClient() (*mongo.Client, context.Context, context.CancelFunc) {
	return utils.GetMongoConnection(u.mongoURL)
}

func (u *TransactionRepositoryImpl) getCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(u.mongoDB).Collection("transaction_data")
}

func (u *TransactionRepositoryImpl) findSingleResult(ctx context.Context, collection *mongo.Collection, cancel context.CancelFunc,
	filter interface{}) (*models.TransactionData, error) {

	result := models.TransactionData{}
	singleResult := collection.FindOne(ctx, filter)
	err := singleResult.Decode(&result)
	defer cancel()
	if err != nil {
		u.logger.Error("Error in decoding transaction data received from Mongo", err)
		return nil, err
	}
	return &result, nil
}

func (u *TransactionRepositoryImpl) findResults(ctx context.Context, collection *mongo.Collection, cancel context.CancelFunc,
	filter interface{}, page int, pageSize int) ([]models.TransactionData, error) {

	var results []models.TransactionData
	allResults, err := collection.Find(ctx, filter, utils.GetFindOptions(page, pageSize))
	defer cancel()
	if err != nil {
		u.logger.Error("Error in finding transaction data from Mongo", err)
		return nil, err
	}
	for allResults.Next(ctx) {
		var singleResult *models.TransactionData
		err := allResults.Decode(&singleResult)
		if err != nil {
			u.logger.Error("Error in decoding transaction data received from Mongo", err)
			return nil, err
		}
		results = append(results, *singleResult)
	}
	if err := allResults.Err(); err != nil {
		u.logger.Error("Error in finding transaction data from Mongo", err)
		return nil, err
	}
	if err := allResults.Close(ctx); err != nil {
		u.logger.Error("Error in closing Mongo connection", err)
		return nil, err
	}
	return results, nil
}
