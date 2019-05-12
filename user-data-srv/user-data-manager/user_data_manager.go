package userdatamanager

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	logger "github.com/jex-lin/golang-logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Loding Logger
var log = logger.NewLogFile(ConfigLogFile)

//Generates A 65 char long user id
func generateRandomStringForID() string {
	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, ConfigUserIDLength)
	for i := range b {
		b[i] = ConfigUserIDCharset[seededRand.Intn(len(ConfigUserIDCharset))]
	}
	return string(b)
}

//UserDataManager The Main User Data Manager Which Will Communicate With MongoDB
type UserDataManager struct {
	Collection *mongo.Collection
	Ctx        context.Context
}

//Init Connects To MongoDB
func (userDataManager *UserDataManager) Init() {
	log.Debug("Initializing ...")
	var err error

	client, _ := mongo.NewClient(options.Client().ApplyURI(ConfigMongoDBAddress))

	//Creating A Timeout Context

	userDataManager.Ctx, _ = context.WithTimeout(context.Background(), ConfigTimeout)

	//Doing The Actual Connection
	err = client.Connect(userDataManager.Ctx)

	if err != nil {
		log.Criticalf("Error Caused At UserDataManager Initialization While Connecting To MongoDB: %s", err)
	}

	//Connection To The Collectiong Which This struct Will Use
	userDataManager.Collection = client.Database(ConfigZeroTechhDB).
		Collection(ConfigUserDataCollection)
}

//generateID Generates A New ID
func (userDataManager UserDataManager) generateID() string {
	log.Debug("Generating New User ID")
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
	log.Debugf("Adding User With Data: %s", user)
	//Checking If Username Or Email Exist
	if userDataManager.doesFieldValueExist(ConfigUsernameField, user[ConfigUsernameField]) {
		return UsernameExistsMsg, nil
	} else if userDataManager.doesFieldValueExist(ConfigEmailField, user[ConfigEmailField]) {
		return EmailExistsMsg, nil
	}
	//Generating A Unique ID
	userID := userDataManager.generateID()
	user[ConfigUserIDField] = userID

	//Adding Some Extra Data
	user[ConfigUserExtraDataField].(map[string]interface{})[ConfigAccountCreationUTCField] = time.Now().Unix()
	user[ConfigUserExtraDataField].(map[string]interface{})[ConfigAccountStatusField] = ConfigAccountStatusActive

	_, err := userDataManager.Collection.InsertOne(userDataManager.Ctx, user)
	return userID, err
}

//GetUserByUsernameOrEmail Returns User Data Based On Username Or Email
func (userDataManager UserDataManager) GetUserByUsernameOrEmail(username, email string, keepPwdOpt ...bool) (map[string]interface{}, error) {
	if username == "" && email == "" {
		return nil, errors.New("No Username And Email Were Passed")
	}

	keepPwd := false
	if len(keepPwdOpt) > 0 {
		keepPwd = keepPwdOpt[0]
	}

	var filter bson.M

	if username != "" {
		filter = bson.M{ConfigUsernameField: username}
	} else if email != "" {
		filter = bson.M{ConfigEmailField: username}
	}

	var user bson.M
	err := userDataManager.Collection.FindOne(userDataManager.Ctx, filter).Decode(&user)
	if user == nil {
		return nil, UserDoesNotExistError{}
	}

	if keepPwd {
		return user, err
	}
	delete(user, ConfigPasswordField) //Removing The Password Field
	return user, err
}

//AuthUser Auths User Based On Username And Password
func (userDataManager UserDataManager) AuthUser(username, email, password string) (bool, error) {
	log.Debugf("Authing User %s With Password %s", username, password)

	//TODO Add Hashing To Check Password
	user, err := userDataManager.GetUserByUsernameOrEmail(username, email, true)
	if err != nil && user == nil {
		return false, errors.New("Invalid Username Or Email")
	} else if err != nil {
		return false, err
	}
	valid := user[ConfigPasswordField] == password
	if !valid {
		return false, errors.New("Invalid Password")
	}
	return valid, nil
}

//GetUser Returns A User Based On UserID
func (userDataManager UserDataManager) GetUser(userID string) (map[string]interface{}, error) {
	log.Debugf("Trying To Get User With UserID: %s", userID)
	var user bson.M
	filter := bson.M{ConfigUserIDField: userID}
	err := userDataManager.Collection.FindOne(userDataManager.Ctx, filter).Decode(&user)
	if user == nil {
		return nil, UserDoesNotExistError{}
	}
	delete(user, ConfigPasswordField) //Removing The Password Field
	return user, err
}

//UpdateUser Updates A Field Of A User
func (userDataManager UserDataManager) UpdateUser(userID, field, newValue string) error {
	log.Debugf("Updating %s Field To %s Of User ID %s", field, newValue, userID)
	filter := bson.M{ConfigUserIDField: userID}
	update := bson.M{"$set": bson.M{field: newValue}}
	_, err := userDataManager.Collection.UpdateOne(userDataManager.Ctx, filter, update)
	return err
}

//DeleteUser Marks User's Account Status As Deleted
func (userDataManager UserDataManager) DeleteUser(userID, username, password string) error {
	valid, _ := userDataManager.AuthUser(username, "", password)
	if valid {
		log.Debugf("Deleting User With UserID %s", userID)
		return userDataManager.UpdateUser(userID, fmt.Sprintf("%s.%s", ConfigUserExtraDataField, ConfigAccountStatusField),
			ConfigAccountStatusDeleted)
	}
	return errors.New("Invalid Auth Info")
}
