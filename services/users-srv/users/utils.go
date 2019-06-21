package users

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/SonicRoshan/Velocity/global/config"
)

func generateUserMetaData() config.UserMeta {
	return config.UserMeta{
		AccountStatus:      config.UserDataConfigAccountStatusActive,
		AccountCreationUTC: time.Now().Unix(),
	}
}

func fieldExist(
	ctx context.Context, filter interface{}, collection *mongo.Collection) (bool, error) {

	var data interface{}
	err := collection.FindOne(ctx, filter).Decode(&data)
	return data != nil, err
}
