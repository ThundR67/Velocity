package verification

import (
	"context"

	"github.com/SonicRoshan/Velocity/global/config"
	logger "github.com/SonicRoshan/Velocity/global/logs"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var log = logger.GetLogger("verification_code_storage.log")

//CodeStore is used to generate, store and verify verification code
type CodeStore struct {
	client         *mongo.Client
	database       *mongo.Database
	mainCollection *mongo.Collection
	contextCancel  context.CancelFunc
	ctx            context.Context
}

//createClient is used to create a client to MongoDB
func (codeStore *CodeStore) createClient() error {
	log.Debug("Creating A Client To MongoDB")
	var err error

	codeStore.client, err = mongo.NewClient(
		options.Client().ApplyURI(config.DBConfigMongoDBAddress),
	)

	if err != nil {
		err = errors.Wrap(err, "Error While Creating Client To MongoDB")
		return err
	}
	log.Info("Created A Client To MongoDB")
	return nil
}

//connect is used to connect to MongoDB
func (codeStore *CodeStore) connect() error {

	log.Debug(
		"Connecting To MongoDB",
		zap.String("DB Name", config.DBConfigMainDB),
	)

	codeStore.ctx = context.TODO()
	err := codeStore.client.Connect(codeStore.ctx)
	if err != nil {
		log.Debug(
			"Connecting To MongoDB Returned Error",
			zap.String("DB Name", config.DBConfigMainDB),
			zap.Error(err),
		)
		err = errors.Wrap(err, "Error While Connecting To MongoDB")
		return err
	}
	codeStore.database = codeStore.client.Database(config.DBConfigMainDB)
	log.Info(
		"Connected To MongoDB",
		zap.String("DB Name", config.DBConfigMainDB),
	)
	return nil
}

//connectToCollections is used to connect to all required collections
func (codeStore *CodeStore) connectToCollections() {
	codeStore.mainCollection = codeStore.database.Collection(config.VerificationMainCollection)
}

//Init is used to initialize codeStore struct
func (codeStore *CodeStore) Init() error {
	err := codeStore.createClient()
	if err != nil {
		return err
	}
	err = codeStore.connect()
	if err != nil {
		return err
	}
	codeStore.connectToCollections()
	return nil
}

//Disconnect is used to disconnect from the mongodb
func (codeStore CodeStore) Disconnect() {
	codeStore.client.Disconnect(codeStore.ctx)
	codeStore.contextCancel()
}
