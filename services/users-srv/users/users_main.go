package users

import (
	"context"

	"github.com/SonicRoshan/Velocity/global/config"
	logger "github.com/SonicRoshan/Velocity/global/logs"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var log = logger.GetLogger("users_low_level_manager.log")
var validateLog = logger.GetLogger("users_data_validater.log")

//Generates A Version 4 UUID
func generateUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		err = errors.Wrap(err, "Error While Generating UUID")
	}
	return id.String(), err
}

//Users is used to handle user data
type Users struct {
	DBName          string
	client          *mongo.Client
	database        *mongo.Database
	mainCollection  *mongo.Collection
	extraCollection *mongo.Collection
	metaCollection  *mongo.Collection
	contextCancel   context.CancelFunc
	ctx             context.Context
}

//createClient is used to create a client to MongoDB
func (users *Users) createClient() {
	log.Debug("Creating A Client To MongoDB")

	users.client, _ = mongo.NewClient(
		options.Client().ApplyURI(config.DBConfigMongoDBAddress),
	)

	log.Info("Created A Client To MongoDB")
}

//connect is used to connect to MongoDB
func (users *Users) connect() {

	log.Debug(
		"Connecting To MongoDB",
		zap.String("DB Name", users.DBName),
	)

	users.ctx = context.TODO()
	users.client.Connect(users.ctx)
	users.database = users.client.Database(users.DBName)

	log.Info(
		"Connected To MongoDB",
		zap.String("DB Name", users.DBName),
	)

}

//connectToCollections is used to connect to all required collections
func (users *Users) connectToCollections() {
	users.mainCollection = users.database.Collection(config.DBConfigUserMainDataCollection)
	users.extraCollection = users.database.Collection(config.DBConfigUserExtraDataCollection)
	users.metaCollection = users.database.Collection(config.DBConfigUserMetaDataCollection)
}

//doesUsernameOrEmailExists is used  to check if username or email exist
func (users Users) doesUsernameOrEmailExists(mainData config.UserMain) (bool, string, error) {

	log.Debug(
		"Checking if username %s or email %s exists",
		zap.String("Email", mainData.Email),
		zap.String("Username", mainData.Username),
	)

	usernameExist, _ := fieldExist(users.ctx,
		config.UserMain{Username: mainData.Username},
		users.mainCollection)

	emailExist, _ := fieldExist(users.ctx,
		config.UserMain{Email: mainData.Email},
		users.mainCollection)

	if usernameExist {
		log.Info("Username Does Exist")
		return true, config.UsernameExistMsg, nil
	} else if emailExist {
		log.Info("Email Does Exist")
		return true, config.EmailExistMsg, nil
	}

	log.Info("Username And Email Dont Exists")
	return false, "", nil
}

//getFilterByUsernameOrEmail is used to get a mongodb filter based on either username or email
func (users Users) getFilterByUsernameOrEmail(username, email string) config.UserMain {
	if username != "" {
		return config.UserMain{Username: username}
	} else if email != "" {
		return config.UserMain{Email: email}
	}
	return config.UserMain{}
}

//generateID is used to generate a v4 UUID
func (users Users) generateID() (string, error) {
	log.Debug("Generating UUID")
	var userID string
	var err error
	userIDExists := true

	for userIDExists {
		userID, err = generateUUID()
		if err != nil {
			log.Error("Generating UUID Returned Error", zap.Error(err))
			err = errors.Wrap(err, "Error While Generating UUID v4")
			return "", err
		}
		userIDExists, _ = fieldExist(users.ctx,
			config.UserMain{UserID: userID},
			users.mainCollection)
	}

	log.Info("Generated UUID", zap.String("UUID", userID))
	return userID, nil
}

//Init is used to initialize users struct
func (users *Users) Init() {
	if users.DBName == "" {
		users.DBName = config.DBConfigMainDB
	}
	users.createClient()
	users.connect()
	users.connectToCollections()
}

//Disconnect is used to disconnect from the mongodb
func (users Users) Disconnect() {
	users.client.Disconnect(users.ctx)
}
