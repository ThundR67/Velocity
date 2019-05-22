package userdatamanager

import (
	"context"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Generates A Version 4 UUID
func generateUUID() (string, error) {
	id, err := uuid.NewRandom()
	return id.String(), err
}

//UserDataManager The Main User Data Manager Which Will Communicate With MongoDB
type UserDataManager struct {
	database                *mongo.Database
	userDataCollection      *mongo.Collection
	userExtraDataCollection *mongo.Collection
	userMetaDataCollection  *mongo.Collection
	ctx                     context.Context
}

//Init Connects To MongoDB
func (userDataManager *UserDataManager) Init(dbname ...string) error {
	dbName := config.DBConfigZeroTechhDB
	if dbname != nil {
		dbName = dbname[0]
	}
	//Creating the client
	client, err := mongo.NewClient(options.Client().ApplyURI(config.DBConfigMongoDBAddress))
	if err != nil {
		err = errors.Wrap(err, "Error While Creating Client To MongoDB")
		return err
	}
	userDataManager.ctx = context.TODO()

	//Doing The Actual Connection
	err = client.Connect(userDataManager.ctx)
	if err != nil {
		err = errors.Wrap(err, "Error While Connecting To MongoDB")
		return err
	}

	//Connection To The DB Which This struct Will Use
	userDataManager.database = client.Database(dbName)

	//Connecting to all required collections
	userDataManager.userDataCollection = userDataManager.database.Collection(config.DBConfigUserDataCollection)
	userDataManager.userExtraDataCollection = userDataManager.database.Collection(config.DBConfigUserExtraDataCollection)
	userDataManager.userMetaDataCollection = userDataManager.database.Collection(config.DBConfigUserMetaDataCollection)

	return err
}

//generateID Generates A New ID
func (userDataManager UserDataManager) generateID() (string, error) {
	var userID string
	var err error
	idFound := false
	userID = "ecbb811d-8be4-446e-b46d-45c1ddf4e606"
	for !idFound {
		userID, err = generateUUID()
		if err != nil {
			err = errors.Wrap(err, "Error While Generating UUID v4")
			return "", err
		}
		/* idFound Will Be True If No User With UserID Exist,
		If Exist Then idFound Will Be True And New ID Will Be Generated */
		user, _ := readDocument(userDataManager.ctx,
			config.DBConfigUserIDField,
			userID,
			userDataManager.userDataCollection)
		idFound = user == nil
	}
	return userID, nil
}

//AddUser Adds An User To DB
func (userDataManager UserDataManager) AddUser(userMainData, userExtraData map[string]interface{}) (string, string, error) {
	//Checking if user data is valid
	if !validateUserMainData(userMainData) || !validateUserExtraData(userExtraData) {
		return "", config.InvalidUserDataMsg, nil
	}

	//Checking If Username Or Email Exist
	usernameExist, _ := fieldExist(userDataManager.ctx, config.DBConfigUsernameField, userMainData[config.DBConfigUsernameField], userDataManager.userDataCollection)
	emailExist, _ := fieldExist(userDataManager.ctx, config.DBConfigEmailField, userMainData[config.DBConfigEmailField], userDataManager.userDataCollection)

	if usernameExist {
		return "", config.UsernameExistMsg, nil
	} else if emailExist {
		return "", config.EmailExistMsg, nil
	}

	//making the meta data
	userMetaData := generateUserMetaData()

	//Generating A Unique ID
	userID, err := userDataManager.generateID()
	if err != nil {
		err = errors.Wrap(err, "Error While Generating ID")
		return "", "", err
	}

	//Adding that id
	userMainData[config.DBConfigUserIDField] = userID
	userExtraData[config.DBConfigUserIDField] = userID
	userMetaData[config.DBConfigUserIDField] = userID

	//Adding Data To Thier Respective Collections
	createDocument(userDataManager.ctx, userMainData, userDataManager.userDataCollection)
	createDocument(userDataManager.ctx, userExtraData, userDataManager.userExtraDataCollection)
	createDocument(userDataManager.ctx, userMetaData, userDataManager.userMetaDataCollection)
	return userID, "", nil
}

//GetUserByUsernameOrEmail Returns User Data Based On Username Or Email
func (userDataManager UserDataManager) GetUserByUsernameOrEmail(username, email string) (map[string]interface{}, string, error) {
	if username == "" && email == "" {
		return nil, config.InvalidUsernameAndEmailMsg, nil
	}

	//Checking of either username or email was given
	var key, value string
	if username != "" {
		key = config.DBConfigUsernameField
		value = username
	} else if email != "" {
		key = config.DBConfigEmailField
		value = email
	}

	//Reading document
	userData, err := readDocument(userDataManager.ctx, key, value, userDataManager.userDataCollection)
	if err != nil {
		err = errors.Wrap(err, "Error While Reading Document")
		return nil, "", err
	} else if userData == nil {
		return nil, config.InvalidUsernameOrEmailMsg, nil
	}

	return userData, "", nil
}

//AuthUser Auths User Based On Username And Password
func (userDataManager UserDataManager) AuthUser(username, email, password string) (bool, string, string, error) {

	//Checking email or username
	user, msg, err := userDataManager.GetUserByUsernameOrEmail(username, email)
	if err != nil || user == nil || msg != "" {
		return false, "", config.InvalidUsernameOrEmailMsg, nil
	}
	//Checking password
	valid := user[config.DBConfigPasswordField] == password
	if !valid {
		return false, "", config.InvalidPasswordMsg, nil
	}

	return valid, user[config.DBConfigUserIDField].(string), "", nil
}

//GetUserData Returns User Data Based On UserID
func (userDataManager UserDataManager) GetUserData(userID, collection string) (map[string]interface{}, error) {
	data, err := readDocument(userDataManager.ctx, config.DBConfigUserIDField, userID, userDataManager.database.Collection((collection)))
	if err != nil {
		err = errors.Wrap(err, "Error While Reading Document")
		return nil, err
	}
	delete(data, config.DBConfigPasswordField) //Removing The Password Field
	return data, err
}

//UpdateUserData Updates A Field Of A User
func (userDataManager UserDataManager) UpdateUserData(userID, field, newValue, collection string) error {
	return updateDocument(userDataManager.ctx, config.DBConfigUserIDField, userID, field, newValue, userDataManager.database.Collection(collection))
}

//DeleteUser Marks User's Account Status As Deleted
func (userDataManager UserDataManager) DeleteUser(userID, username, password string) (string, error) {
	valid, _, _, _ := userDataManager.AuthUser(username, "", password)
	if valid {
		return "", updateDocument(userDataManager.ctx, config.DBConfigUserIDField, userID, config.DBConfigAccountStatusField, config.UserDataConfigAccountStatusDeleted, userDataManager.userMetaDataCollection)
	}
	return config.InvalidAuthDataMsg, nil
}
