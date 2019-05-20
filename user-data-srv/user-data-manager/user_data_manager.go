package userdatamanager

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/SonicRoshan/Velocity/global/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Generates A 65 char long user id
func generateRandomStringForID() string {
	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, config.UserDataConfigUserIDLength)
	for i := range b {
		b[i] = config.UserDataConfigUserIDCharset[seededRand.Intn(len(config.UserDataConfigUserIDCharset))]
	}
	return string(b)
}

//UserDataManager The Main User Data Manager Which Will Communicate With MongoDB
type UserDataManager struct {
	Collection *mongo.Collection
	Ctx        context.Context
}

//Init Connects To MongoDB
func (userDataManager *UserDataManager) Init() error {

	client, err := mongo.NewClient(options.Client().ApplyURI(config.DBConfigMongoDBAddress))
	if err != nil {
		return err
	}

	//Creating A Timeout Context

	userDataManager.Ctx = context.TODO()

	//Doing The Actual Connection
	err = client.Connect(userDataManager.Ctx)
	if err != nil {
		return err
	}

	//Connection To The Collectiong Which This struct Will Use
	userDataManager.Collection = client.Database(config.DBConfigZeroTechhDB).
		Collection(config.DBConfigUserDataCollection)
	return err
}

//generateID Generates A New ID
func (userDataManager UserDataManager) generateID() string {
	var userID string
	idFound := false
	for !idFound {
		userID = generateRandomStringForID()
		/* idFound Will Be True If No User With UserID Exist,
		If Exist Then idFound Will Be True And New ID Will Be Generated */
		user, _ := userDataManager.GetUser(userID)
		idFound = user == nil
	}
	return userID
}

//doesFieldValueExist Checks If A Users Particular Field Has An Value
func (userDataManager UserDataManager) doesFieldValueExist(field string, value interface{}) bool {
	var user bson.M
	filter := bson.M{field: value}
	userDataManager.Collection.FindOne(userDataManager.Ctx, filter).Decode(&user)
	return user != nil
}

//AddUser Adds An User To DB
func (userDataManager UserDataManager) AddUser(user map[string]interface{}) (string, error) {
	//Checking if user data is valid
	if !validateUserData(user) {
		return "", config.InvalidUserDataError
	}

	//Checking If Username Or Email Exist
	if userDataManager.doesFieldValueExist(config.DBConfigUsernameField, user[config.DBConfigUsernameField]) {
		return "", config.UsernameExistError
	} else if userDataManager.doesFieldValueExist(config.DBConfigEmailField, user[config.DBConfigEmailField]) {
		return "", config.EmailExistError
	}

	//Generating A Unique ID
	userID := userDataManager.generateID()
	user[config.DBConfigUserIDField] = userID

	//Adding Some Extra Data
	user[config.DBConfigUserExtraDataField].(map[string]interface{})[config.DBConfigAccountCreationUTCField] = time.Now().Unix()
	user[config.DBConfigUserExtraDataField].(map[string]interface{})[config.DBConfigAccountStatusField] = config.UserDataConfigAccountStatusActive

	_, err := userDataManager.Collection.InsertOne(userDataManager.Ctx, user)
	return userID, err
}

//GetUserByUsernameOrEmail Returns User Data Based On Username Or Email
func (userDataManager UserDataManager) GetUserByUsernameOrEmail(username, email string, keepPwdOpt ...bool) (map[string]interface{}, error) {
	if username == "" && email == "" {
		return nil, config.InvalidUsernameAndEmailError
	}

	keepPwd := false
	if len(keepPwdOpt) > 0 {
		keepPwd = keepPwdOpt[0]
	}

	var filter bson.M

	if username != "" {
		filter = bson.M{config.DBConfigUsernameField: username}
	} else if email != "" {
		filter = bson.M{config.DBConfigEmailField: username}
	}

	var user bson.M
	err := userDataManager.Collection.FindOne(userDataManager.Ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, config.UserDoesNotExistError
	}

	if keepPwd {
		return user, err
	}
	delete(user, config.DBConfigPasswordField) //Removing The Password Field
	return user, err
}

//AuthUser Auths User Based On Username And Password
func (userDataManager UserDataManager) AuthUser(username, email, password string) (bool, string, error) {

	//TODO Add Hashing To Check Password
	user, err := userDataManager.GetUserByUsernameOrEmail(username, email, true)
	if err != nil && user == nil {
		return false, "", config.InvalidUsernameOrEmailError
	} else if err != nil {
		return false, "", err
	}
	valid := user[config.DBConfigPasswordField] == password
	if !valid {
		return false, "", config.InvalidPasswordError
	}

	return valid, user[config.DBConfigUserIDField].(string), nil
}

//GetUser Returns A User Based On UserID
func (userDataManager UserDataManager) GetUser(userID string) (map[string]interface{}, error) {
	var user bson.M
	filter := bson.M{config.DBConfigUserIDField: userID}
	err := userDataManager.Collection.FindOne(userDataManager.Ctx, filter).Decode(&user)
	if user == nil {
		return nil, config.UserDoesNotExistError
	}
	delete(user, config.DBConfigPasswordField) //Removing The Password Field
	return user, err
}

//UpdateUser Updates A Field Of A User
func (userDataManager UserDataManager) UpdateUser(userID, field, newValue string) error {
	filter := bson.M{config.DBConfigUserIDField: userID}
	update := bson.M{"$set": bson.M{field: newValue}}
	_, err := userDataManager.Collection.UpdateOne(userDataManager.Ctx, filter, update)
	return err
}

//DeleteUser Marks User's Account Status As Deleted
func (userDataManager UserDataManager) DeleteUser(userID, username, password string) error {
	valid, _, _ := userDataManager.AuthUser(username, "", password)
	if valid {
		return userDataManager.UpdateUser(userID, fmt.Sprintf("%s.%s", config.DBConfigUserExtraDataField, config.DBConfigAccountStatusField),
			config.UserDataConfigAccountStatusDeleted)
	}
	return config.InvalidAuthDataError
}
