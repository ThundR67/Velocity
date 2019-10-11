package verification

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/SonicRoshan/Velocity/global/logger"
	"github.com/SonicRoshan/Velocity/global/utils"
)

var log = logger.GetLogger("verification_code_storage.log")

//CodeStore is used to generate, store and verify verification code
type CodeStore struct {
	client         *mongo.Client
	database       *mongo.Database
	mainCollection *mongo.Collection
}

//Init is used to initialize codeStore struct
func (codeStore *CodeStore) Init() {
	codeStore.client = utils.CreateMongoDB(config.DBConfigMongoDBAddress, log)
	codeStore.database = codeStore.client.Database(config.DBConfigMainDB)
	codeStore.mainCollection = codeStore.database.Collection(
		config.VerificationMainCollection,
	)
}

//Disconnect is used to disconnect from the mongodb
func (codeStore CodeStore) Disconnect() {
	codeStore.client.Disconnect(context.TODO())
}
