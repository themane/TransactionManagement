package repositories

import (
	constants "TxnManagement/contants"
	"TxnManagement/repositories/models"
	"TxnManagement/repositories/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AdminRepositoryImpl struct {
	mongoURL string
	mongoDB  string
	logger   *constants.LoggingUtils
}

func NewAdminRepository(mongoURL string, mongoDB string, logLevel string) *AdminRepositoryImpl {
	return &AdminRepositoryImpl{
		mongoURL: mongoURL,
		mongoDB:  mongoDB,
		logger:   constants.NewLoggingUtils("ADMIN_REPOSITORY", logLevel),
	}
}

func (u *AdminRepositoryImpl) FindById(id string) (*models.AdminData, error) {
	client, ctx := u.getMongoClient()
	defer utils.Disconnect(client, ctx)
	filter := bson.M{"_id": id}
	return u.findSingleResult(ctx, u.getCollection(client), filter)
}

func (u *AdminRepositoryImpl) FindByEmail(userId string) (*models.AdminData, error) {
	client, ctx := u.getMongoClient()
	defer utils.Disconnect(client, ctx)
	filter := bson.M{"email": userId}
	return u.findSingleResult(ctx, u.getCollection(client), filter)
}

func (u *AdminRepositoryImpl) AddUser(adminData models.AdminData) error {
	client, ctx := u.getMongoClient()
	defer utils.Disconnect(client, ctx)
	_, err := u.getCollection(client).InsertOne(ctx, adminData)
	if err != nil {
		u.logger.Error("Error in saving admin data to Mongo", err)
		return err
	}
	return nil
}

func (u *AdminRepositoryImpl) getMongoClient() (*mongo.Client, context.Context) {
	return utils.GetMongoConnection(u.mongoURL)
}

func (u *AdminRepositoryImpl) getCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(u.mongoDB).Collection("admin_data")
}

func (u *AdminRepositoryImpl) findSingleResult(ctx context.Context, collection *mongo.Collection,
	filter interface{}) (*models.AdminData, error) {

	result := models.AdminData{}
	singleResult := collection.FindOne(ctx, filter)
	err := singleResult.Decode(&result)
	if err != nil {
		u.logger.Error("Error in decoding admin data received from Mongo", err)
		return nil, err
	}
	return &result, nil
}
