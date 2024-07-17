package repositories

import (
	constants "TxnManagement/contants"
	"TxnManagement/repositories/models"
	"TxnManagement/repositories/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CustomerRepositoryImpl struct {
	mongoURL string
	mongoDB  string
	logger   *constants.LoggingUtils
}

func NewCustomerRepository(mongoURL string, mongoDB string, logLevel string) *CustomerRepositoryImpl {
	return &CustomerRepositoryImpl{
		mongoURL: mongoURL,
		mongoDB:  mongoDB,
		logger:   constants.NewLoggingUtils("CUSTOMER_REPOSITORY", logLevel),
	}
}

func (u *CustomerRepositoryImpl) FindByIds(ids []string) ([]models.CustomerData, error) {
	client, ctx, cancel := u.getMongoClient()
	defer utils.Disconnect(client, ctx)
	filter := bson.M{"_id": bson.M{"$in": ids}}
	return u.findResults(ctx, u.getCollection(client), cancel, filter, options.Find())
}

func (u *CustomerRepositoryImpl) FindByName(name string, page int, pageSize int) ([]models.CustomerData, error) {
	client, ctx, cancel := u.getMongoClient()
	defer utils.Disconnect(client, ctx)
	filter := bson.M{"name": name}
	return u.findResults(ctx, u.getCollection(client), cancel, filter, utils.GetFindOptions(page, pageSize))
}

func (u *CustomerRepositoryImpl) FindByMobile(mobile string) (*models.CustomerData, error) {
	client, ctx, cancel := u.getMongoClient()
	defer utils.Disconnect(client, ctx)
	filter := bson.M{"mobiles": mobile}
	return u.findSingleResult(ctx, u.getCollection(client), cancel, filter)
}

func (u *CustomerRepositoryImpl) FindByVillage(village string, page int, pageSize int) ([]models.CustomerData, error) {
	client, ctx, cancel := u.getMongoClient()
	defer utils.Disconnect(client, ctx)
	filter := bson.M{"address.village": village}
	return u.findResults(ctx, u.getCollection(client), cancel, filter, utils.GetFindOptions(page, pageSize))
}

func (u *CustomerRepositoryImpl) FindByKeyword(keyword string, page int, pageSize int) ([]models.CustomerData, error) {
	client, ctx, cancel := u.getMongoClient()
	defer utils.Disconnect(client, ctx)
	filter := bson.M{
		"$or": bson.A{
			bson.M{"name": bson.M{"$regex": keyword}},
			bson.M{"father_name": bson.M{"$regex": keyword}},
			bson.M{"mobiles": bson.M{"$regex": keyword}},
			bson.M{"address.village": bson.M{"$regex": keyword}},
		},
	}
	return u.findResults(ctx, u.getCollection(client), cancel, filter, utils.GetFindOptions(page, pageSize))
}

func (u *CustomerRepositoryImpl) AddUser(customerData models.CustomerData) error {
	client, ctx, cancel := u.getMongoClient()
	defer utils.Disconnect(client, ctx)
	_, err := u.getCollection(client).InsertOne(ctx, customerData)
	defer cancel()
	if err != nil {
		u.logger.Error("Error in saving customer data to Mongo", err)
		return err
	}
	return nil
}

func (u *CustomerRepositoryImpl) getMongoClient() (*mongo.Client, context.Context, context.CancelFunc) {
	return utils.GetMongoConnection(u.mongoURL)
}

func (u *CustomerRepositoryImpl) getCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(u.mongoDB).Collection("customer_data")
}

func (u *CustomerRepositoryImpl) findSingleResult(ctx context.Context, collection *mongo.Collection, cancel context.CancelFunc,
	filter interface{}) (*models.CustomerData, error) {

	result := models.CustomerData{}
	singleResult := collection.FindOne(ctx, filter)
	defer cancel()
	err := singleResult.Decode(&result)
	if err != nil {
		u.logger.Error("Error in decoding customer data received from Mongo", err)
		return nil, err
	}
	return &result, nil
}

func (u *CustomerRepositoryImpl) findResults(ctx context.Context, collection *mongo.Collection, cancel context.CancelFunc,
	filter interface{}, options *options.FindOptions) ([]models.CustomerData, error) {

	var results []models.CustomerData
	allResults, err := collection.Find(ctx, filter, options)
	defer cancel()
	if err != nil {
		u.logger.Error("Error in finding customer data from Mongo", err)
		return nil, err
	}
	for allResults.Next(ctx) {
		var singleResult *models.CustomerData
		err := allResults.Decode(&singleResult)
		if err != nil {
			u.logger.Error("Error in decoding customer data received from Mongo", err)
			return nil, err
		}
		results = append(results, *singleResult)
	}
	if err := allResults.Err(); err != nil {
		u.logger.Error("Error in finding customer data from Mongo", err)
		return nil, err
	}
	if err := allResults.Close(ctx); err != nil {
		u.logger.Error("Error in closing Mongo connection", err)
		return nil, err
	}
	return results, nil
}
