package utils

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetFindOptions(page int, pageSize int) *options.FindOptions {
	findOptions := options.Find()
	findOptions.SetSkip(int64((page - 1) * pageSize))
	findOptions.SetLimit(int64(pageSize))
	return findOptions
}
