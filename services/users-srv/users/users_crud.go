package users

import (
	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//GetByUsernameOrEmail is used to get user based on username or email
func (users Users) GetByUsernameOrEmail(username, email string) (config.UserMain, string) {

	log.Debug(
		"Getting User By Username Or Email",
		zap.String("Username", username),
		zap.String("Email", email),
	)

	filter := users.getFilterByUsernameOrEmail(username, email)
	var user config.UserMain
	err := users.mainCollection.FindOne(users.ctx, filter).Decode(&user)

	if err != nil {
		return config.UserMain{}, config.InvalidUsernameOrEmailMsg
	}

	log.Info(
		"Got User By Username Or Email",
		zap.String("Username", username),
		zap.String("Email", email),
		zap.Any("User", user),
	)

	return user, ""
}

//Auth is used to authenticate a user
func (users Users) Auth(username, email, password string) (bool, string, string) {

	log.Debug(
		"Authenticating User",
		zap.String("Username", username),
		zap.String("Email", email),
		zap.String("password", password),
	)

	user, msg := users.GetByUsernameOrEmail(username, email)

	if msg != "" || user == (config.UserMain{}) {
		return false, "", msg
	}

	//Checking password (TODO Add hashing)
	valid := user.Password == password

	log.Info(
		"Authenticated User",
		zap.String("Username", username),
		zap.String("Email", email),
		zap.String("password", password),
	)

	if !valid {
		return false, "", config.InvalidPasswordMsg
	}
	return valid, user.UserID, ""
}

//Add is used to add user to database
func (users Users) Add(
	mainData config.UserMain, extraData config.UserExtra) (string, string) {

	log.Debug(
		"Adding User",
		zap.Any("Main Data", mainData),
		zap.Any("Extra Data", extraData),
	)

	if !isValid(mainData, extraData) {
		log.Info("Invalid User Data")
		return "", config.InvalidUserDataMsg
	}

	exists, msg := users.doesUsernameOrEmailExists(mainData)
	if exists {
		return "", msg
	}

	metaData := generateUserMetaData()
	userID, err := users.generateID()
	if err != nil {
		log.Error(
			"Adding User Returned Error",
			zap.Any("Main Data", mainData),
			zap.Any("Extra Data", extraData),
			zap.Any("Meta Data", metaData),
			zap.Error(err),
		)
		err = errors.Wrap(err, "Error While Adding User")
		return "", ""
	}

	mainData.UserID = userID
	extraData.UserID = userID
	metaData.UserID = userID

	users.mainCollection.InsertOne(users.ctx, mainData)
	users.extraCollection.InsertOne(users.ctx, extraData)
	users.metaCollection.InsertOne(users.ctx, metaData)

	if config.DebugMode {
		users.Activate(mainData.Email)
	}

	log.Info(
		"Added User",
		zap.Any("Main Data", mainData),
		zap.Any("Extra Data", extraData),
		zap.Any("Meta Data", metaData),
	)

	return userID, ""
}

//Get is used to get user's data
func (users Users) Get(userID, collectionName string, data interface{}) error {

	log.Debug(
		"Getting User's Data",
		zap.String("UserID", userID),
		zap.String("Collection", collectionName),
	)

	collection := users.database.Collection(collectionName)
	err := collection.FindOne(users.ctx, config.UserMain{UserID: userID}).Decode(data)

	if err != nil {
		log.Error(
			"Getting User's Data Returned Error",
			zap.String("UserID", userID),
			zap.String("Collection", collectionName),
			zap.Error(err),
		)
		err = errors.Wrap(err, "Error While Getting User Data")
		return err
	}

	log.Info(
		"Got User's Data",
		zap.String("UserID", userID),
		zap.String("Collection", collectionName),
	)

	return nil
}

//Update is used to update users data in any collection
func (users Users) Update(userID string, update interface{}, collectionName string) {

	log.Debug(
		"Updating User's Data",
		zap.String("UserID", userID),
		zap.String("Collection", collectionName),
		zap.Any("Update", update),
	)

	//Todo Validate The Update

	collection := users.database.Collection(collectionName)

	collection.UpdateOne(users.ctx,
		config.UserMain{UserID: userID}, map[string]interface{}{"$set": update})

	log.Info(
		"Updated User's Data",
		zap.String("UserID", userID),
		zap.String("Collection", collectionName),
		zap.Any("Update", update),
	)
}

//Delete is used to mark a user as deleted
func (users Users) Delete(userID, username, password string) string {

	log.Debug(
		"Deleting User",
		zap.String("UserID", userID),
		zap.String("Username", username),
		zap.String("Password", password),
	)

	valid, _, _ := users.Auth(username, "", password)

	update := config.UserMeta{AccountStatus: config.UserDataConfigAccountStatusDeleted}

	if !valid {
		log.Info(
			"Invalid Auth Info",
			zap.String("UserID", userID),
			zap.String("Username", username),
			zap.String("Password", password),
		)
		return config.InvalidAuthDataMsg
	}

	users.Update(userID, update, config.DBConfigUserMetaDataCollection)

	log.Info(
		"Deleted User",
		zap.String("UserID", userID),
		zap.String("Username", username),
		zap.String("Password", password),
	)

	return ""
}

//Activate is used to mark an account as active
func (users Users) Activate(email string) string {

	log.Debug(
		"Activating User With Email",
		zap.String("Email", email),
	)

	userData, msg := users.GetByUsernameOrEmail("", email)
	if msg != "" {
		return msg
	}
	userID := userData.UserID

	log.Debug(
		"Activating User With ID",
		zap.String("UserID", userID),
	)

	update := config.UserMeta{AccountStatus: config.UserDataConfigAccountStatusActive}
	users.Update(userID, update, config.DBConfigUserMetaDataCollection)

	log.Info(
		"Activated User",
		zap.String("UserID", userID),
		zap.String("Email", email),
	)

	return ""
}
