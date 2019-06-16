package users

import (
	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/pkg/errors"
)

//GetByUsernameOrEmail is used to get user based on username or email
func (users Users) GetByUsernameOrEmail(username, email string) (config.UserMain, string) {
	log.Debugf("Getting User By Username %s Email %s", username, email)
	filter := users.getFilterByUsernameOrEmail(username, email)
	var user config.UserMain
	err := users.mainCollection.FindOne(users.ctx, filter).Decode(&user)

	if err != nil {
		return config.UserMain{}, config.InvalidUsernameOrEmailMsg
	}

	log.Infof("Got User By Username %s Email %s", username, email)
	return user, ""
}

//Auth is used to authenticate a user
func (users Users) Auth(username, email, password string) (bool, string, string, error) {
	log.Debugf("Authenticating With Username %s Email %s", username, email)

	user, msg := users.GetByUsernameOrEmail(username, email)

	if msg != "" || user == (config.UserMain{}) {
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
func (users Users) Add(
	mainData config.UserMain, extraData config.UserExtra) (string, string, error) {

	log.Debugf(`Adding User With 
				Main Data %+v
				Extra Data %+v`, mainData, extraData)

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
		log.Errorf(`Adding User With
					Main Data %+v
					User Extra Data %+v
					User Meta Data %+v
					Returned Error %+v`,
			mainData, extraData, metaData, err)
		err = errors.Wrap(err, "Error While Adding User")
		return "", "", err
	}

	mainData.UserID = userID
	extraData.UserID = userID
	metaData.UserID = userID

	users.mainCollection.InsertOne(users.ctx, mainData)
	users.extraCollection.InsertOne(users.ctx, extraData)
	users.metaCollection.InsertOne(users.ctx, metaData)
	log.Infof(`Added User With
					Main Data %+v
					User Extra Data %+v
					User Meta Data %+v`,
		mainData, extraData, metaData)
	return userID, "", nil
}

//Get is used to decoded user data into data interface
func (users Users) Get(userID, collectionName string, data interface{}) error {
	log.Debugf("Getting User Data In Collection %s With UserID %s", collectionName, userID)

	collection := users.database.Collection(collectionName)
	err := collection.FindOne(users.ctx, config.UserMain{UserID: userID}).Decode(data)

	if err != nil {
		log.Errorf("Getting User In Collection %s With UserID %s Returned Error %+v",
			collectionName, userID, err)
		err = errors.Wrap(err, "Error While Getting User Data")
		return err
	}

	log.Infof("Got Users Data In Collection %s With UserID %s", collectionName, userID)
	return nil
}

//Update is used to update users data in any collection
func (users Users) Update(userID string, update interface{}, collectionName string) error {
	log.Debugf("Updating User %s Data In Collection %s With Update %+v",
		userID, collectionName, update)

	collection := users.database.Collection(collectionName)

	_, err := collection.UpdateOne(users.ctx,
		config.UserMain{UserID: userID}, map[string]interface{}{"$set": update})

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

	update := config.UserMeta{AccountStatus: config.UserDataConfigAccountStatusDeleted}

	if !valid {
		log.Infof("Invalid Auth Data While Deleting Username %s", username)
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
