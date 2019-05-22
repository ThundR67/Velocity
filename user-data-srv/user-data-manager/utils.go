package userdatamanager

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/SonicRoshan/Velocity/global/config"
)

func generateUserMetaData() map[string]interface{} {
	return map[string]interface{}{
		config.DBConfigAccountStatusField:      config.UserDataConfigAccountStatusActive,
		config.DBConfigAccountCreationUTCField: time.Now().Unix(),
	}
}

func createDocument(ctx context.Context, data map[string]interface{}, collection *mongo.Collection) {
	collection.InsertOne(ctx, data)
}

func readDocument(ctx context.Context, key string, value interface{}, collection *mongo.Collection) (map[string]interface{}, error) {
	var data map[string]interface{}
	filter := bson.M{key: value}
	err := collection.FindOne(ctx, filter).Decode(&data)
	return data, err
}

func fieldExist(ctx context.Context, key string, value interface{}, collection *mongo.Collection) (bool, error) {
	data, err := readDocument(ctx, key, value, collection)
	return data != nil, err
}

func updateDocument(ctx context.Context, key, value, field, newValue string, collection *mongo.Collection) error {
	filter := bson.M{key: value}
	update := bson.M{"$set": bson.M{field: newValue}}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
