package users

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/SonicRoshan/Velocity/global/config"
)

func generateUserMetaData() config.UserType {
	return config.UserType{
		AccountStatus:      config.UserDataConfigAccountStatusActive,
		AccountCreationUTC: time.Now().Unix(),
	}
}

func fieldExist(ctx context.Context, filter config.UserType, collection *mongo.Collection) (bool, error) {
	var data config.UserType
	err := collection.FindOne(ctx, filter).Decode(&data)
	return data != (config.UserType{}), err
}
