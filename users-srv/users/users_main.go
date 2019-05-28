package users

import (
	"context"

	"github.com/SonicRoshan/Velocity/global/config"
	logger "github.com/SonicRoshan/Velocity/global/logs"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var log = logger.GetLogger("users_low_level_manager.log")

//Generates A Version 4 UUID
func generateUUID() (string, error) {
	id, err := uuid.NewRandom()
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
func (users *Users) createClient() error {
	log.Debug("Creating A Client To MongoDB")
	var err error

	users.client, err = mongo.NewClient(
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
func (users *Users) connect() error {
	log.Debugf("Connecting To MongoDB With DB Name %s", users.DBName)
	users.ctx = context.TODO()
	err := users.client.Connect(users.ctx)
	if err != nil {
		err = errors.Wrap(err, "Error While Connecting To MongoDB")
		return err
	}
	users.database = users.client.Database(users.DBName)
	log.Info("Connected To MongoDB And DB")
	return nil
}

//connectToCollections is used to connect to all required collections
func (users *Users) connectToCollections() {
	users.mainCollection = users.database.Collection(config.DBConfigUserMainDataCollection)
	users.extraCollection = users.database.Collection(config.DBConfigUserExtraDataCollection)
	users.metaCollection = users.database.Collection(config.DBConfigUserMetaDataCollection)
}

//doesUsernameOrEmailExists is used  to check if username or email exist
func (users Users) doesUsernameOrEmailExists(mainData config.UserType) (bool, string, error) {
	log.Debugf("Checking if username %s or email %s exists",
		mainData.Email,
		mainData.Username)

	usernameExist, _ := fieldExist(users.ctx,
		config.UserType{Username: mainData.Username},
		users.mainCollection)

	emailExist, _ := fieldExist(users.ctx,
		config.UserType{Email: mainData.Email},
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
func (users Users) getFilterByUsernameOrEmail(username, email string) config.UserType {
	if username != "" {
		return config.UserType{Username: username}
	} else if email != "" {
		return config.UserType{Username: username}
	}
	return config.UserType{}
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
			log.Errorf("Generating UUID Returned Error %+v", err)
			err = errors.Wrap(err, "Error While Generating UUID v4")
			return "", err
		}
		userIDExists, _ = fieldExist(users.ctx,
			config.UserType{UserID: userID},
			users.mainCollection)
	}
	log.Infof("Generated UUID %s", userID)
	return userID, nil
}

//Init is used to initialize users struct
func (users *Users) Init() error {
	if users.DBName == "" {
		users.DBName = config.DBConfigMainDB
	}
	err := users.createClient()
	if err != nil {
		return err
	}
	err = users.connect()
	if err != nil {
		return err
	}
	users.connectToCollections()
	return nil
}

//Disconnect is used to disconnect from the mongodb
func (users Users) Disconnect() {
	users.client.Disconnect(users.ctx)
	users.contextCancel()
}
