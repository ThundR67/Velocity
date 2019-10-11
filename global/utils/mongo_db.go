package utils

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

//CreateMongoDB creates a mongoDB client
func CreateMongoDB(address string, log *zap.Logger) *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(
			"Can Not Connect To MongoDB At Address",
			zap.String("Address", address),
			zap.Error(err),
		)
		panic("Can Not Connect To MongoDB At Address " + address)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(
			"Can Not Maintain Connection To MongoDB At Address",
			zap.String("Address", address),
			zap.Error(err),
		)
		panic("Can Not Maintain Connection To MongoDB At Address " + address)
	}

	return client
}
