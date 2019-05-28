package users

import (
	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/pkg/errors"
)

//GetByUsernameOrEmail is used to get user based on username or email
func (users Users) GetByUsernameOrEmail(username, email string) (config.UserType, string) {
	log.Debugf("Getting User By Username %s Email %s", username, email)
	filter := users.getFilterByUsernameOrEmail(username, email)
	var user config.UserType
	err := users.mainCollection.FindOne(users.ctx, filter).Decode(&user)

	if err != nil {
		return config.UserType{}, config.InvalidUsernameOrEmailMsg
	}

	log.Infof("Got User By Username %s Email %s", username, email)
	return user, ""
}

//Auth is used to authenticate a user
func (users Users) Auth(username, email, password string) (bool, string, string, error) {
	log.Debugf("Authenticating With Username %s Email %s", username, email)

	user, msg := users.GetByUsernameOrEmail(username, email)

	if msg != "" || user == (config.UserType{}) {
		return false, "", msg, nil
	}

	//Checking password (TODO Add hashing)
	valid := user.Password == password

	log.Infof("Authenticated With Username %s Email %s", username, email)

	if !valid {
		return false, "", config.InvalidPasswordMsg, nil
	}
	return valid, user.UserID, "", nil
}

//Add is used to add user to database
func (users Users) Add(mainData, extraData config.UserType) (string, string, error) {
	log.Debugf("Adding User With Main Data %+v\n Extra %+v", mainData, extraData)
	if !isValid(mainData, extraData) {
		log.Info("Invalid User Data")
		return "", config.InvalidUserDataMsg, nil
	}

	exists, msg, err := users.doesUsernameOrEmailExists(mainData)
	if exists {
		return "", msg, err
	}

	metaData := generateUserMetaData()
	userID, err := users.generateID()
	if err != nil {
		log.Errorf("Adding User Data %+v User Extra Data %+v User Meta Data %+v Returned Error %+v",
			mainData, extraData, metaData, err)
		err = errors.Wrap(err, "Error While Generating ID")
		return "", "", err
	}

	mainData.UserID = userID
	extraData.UserID = userID
	metaData.UserID = userID

	users.mainCollection.InsertOne(users.ctx, mainData)
	users.extraCollection.InsertOne(users.ctx, extraData)
	users.metaCollection.InsertOne(users.ctx, metaData)
	log.Debugf("Added User With Main Data %+v\n Extra %+v", mainData, extraData)
	return userID, "", nil
}

//Get is used to get user data in any of the collection based on user id
func (users Users) Get(userID, collectionName string) (config.UserType, error) {
	log.Debugf("Getting User Data In Collection %s With UserID %v", collectionName, userID)

	var data config.UserType
	collection := users.database.Collection(collectionName)
	err := collection.FindOne(users.ctx, config.UserType{UserID: userID}).Decode(&data)

	if err != nil || data == (config.UserType{}) {
		log.Errorf("Getting User In Collection %s With UserID %s Returned Error %+v",
			collectionName, userID, err)
		err = errors.Wrap(err, "Error While Getting User Data")
		return config.UserType{}, err
	}

	log.Infof("Got Users Data In Collection %s", collectionName)
	return data, nil
}

//Update is used to update users data in any collection
func (users Users) Update(userID string, update config.UserType, collectionName string) error {
	log.Debugf("Updating User %s Data In Collection %s With Update %+v",
		userID, collectionName, update)

	collection := users.database.Collection(collectionName)

	_, err := collection.UpdateOne(users.ctx,
		config.UserType{UserID: userID}, map[string]config.UserType{"$set": update})

	if err != nil {
		log.Errorf("Updating User %s Data In Collection %s With Update %+v Returned Error %+v",
			userID, collectionName, update, err)
		err = errors.Wrap(err, "Error While Updating User Data")
		return err
	}
	log.Infof("Updated User %s Data In Collection %s",
		userID, collectionName)
	return nil
}

//Delete is used to mark a user as deleted
func (users Users) Delete(userID, username, password string) (string, error) {
	log.Debugf("Deleting User %s", username)

	valid, _, _, _ := users.Auth(username, "", password)

	update := config.UserType{AccountStatus: config.UserDataConfigAccountStatusDeleted}

	if !valid {
		log.Infof("Invaling Auth Data While Deleting User %s", username)
		return config.InvalidAuthDataMsg, nil
	}

	err := users.Update(userID, update, config.DBConfigUserMetaDataCollection)
	if err != nil {
		log.Errorf("Deleting User %s Returned Error %+v", username, err)
		err = errors.Wrap(err, "Error While Deleting User")
		return "", err
	}
	log.Infof("Deleted User %s", username)
	return "", err
}
