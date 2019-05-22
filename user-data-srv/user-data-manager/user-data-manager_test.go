package userdatamanager

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/assert"
)

func generateRandomString(length int) string {
	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = config.UserDataConfigUserIDCharset[seededRand.Intn(len(config.UserDataConfigUserIDCharset))]
	}
	return string(b)
}

//Creating dataManager With collection Created To TestDB
func setup() UserDataManager {
	output := UserDataManager{}
	output.Init("TestDB")
	return output
}

var dataManager = setup()

//TestLowLevelUserDataManager Will Test Low Level user data manager
func TestLowLevelUserDataManager(t *testing.T) {
	//Creating mock Data For Testing
	assert := assert.New(t)

	mockUsername := strings.ToLower(generateRandomString(10))
	mockPassword := generateRandomString(10)

	mockUserData := map[string]interface{}{
		config.DBConfigUsernameField: mockUsername,
		config.DBConfigPasswordField: mockPassword,
		config.DBConfigEmailField:    mockUsername,
	}

	mockUserExtraData := map[string]interface{}{
		config.DBConfigGenderField:      "male",
		config.DBConfigFirstNameField:   mockUsername,
		config.DBConfigLastNameField:    mockUsername,
		config.DBConfigBirthdayUTCField: int64(0),
	}

	//Testing AddUser
	userID, _, err := dataManager.AddUser(mockUserData, mockUserExtraData)
	assert.NoError(err, "AddUser Returned Error")

	//testing of all three collection got the data
	data, err := dataManager.GetUserData(userID, config.DBConfigUserDataCollection)
	assert.True(data != nil && err == nil, "AddUser Failed To Add Data To Main User Data Collection")
	data, err = dataManager.GetUserData(userID, config.DBConfigUserExtraDataCollection)
	assert.True(data != nil && err == nil, "AddUser Failed To Add Data To Main User Extra Data Collection")
	data, err = dataManager.GetUserData(userID, config.DBConfigUserMetaDataCollection)
	assert.True(data != nil && err == nil, "AddUser Failed To Add Data To Main User Meta Data Collection")

	//Testing GetUser
	user, err := dataManager.GetUserData(userID, config.DBConfigUserDataCollection)
	assert.NoError(err, "GetUserData Returned Error")
	assert.Equal(mockUsername, user[config.DBConfigUsernameField], "Failed At GetUser, Mock Username And Username In DB Must Be Equal")

	//Testing GetUserByUsernameOrEmail
	user, _, err = dataManager.GetUserByUsernameOrEmail(mockUsername, "")
	assert.NoError(err, "GetUserByUsernameOrEmail Returned Error")
	assert.Equal(mockUsername, user[config.DBConfigUsernameField], "Failed At GetUserByUsernameOrEmail, Mock Username And Username In DB Must Be Equal")

	//Testing auth user
	valid, _, _, err := dataManager.AuthUser(mockUsername, "", mockPassword)
	assert.NoError(err, "AuthUser Returned Error")
	assert.True(valid, "Failed At AuthUser For Checking With Correct Password")
	valid, _, _, err = dataManager.AuthUser(mockUsername, "", "wrong-password")
	assert.False(valid, "Failed At AuthUser For Checking With Incorrect Password")

	//Testing UpdateUser
	newPassword := generateRandomString(10)
	dataManager.UpdateUserData(userID, config.DBConfigPasswordField, newPassword, config.DBConfigUserDataCollection)
	user, _, err = dataManager.GetUserByUsernameOrEmail(mockUsername, "")
	assert.NoError(err, "GetUserByUsernameOrEmail Returned Error")
	assert.Equal(newPassword, user[config.DBConfigPasswordField], "Failed At UpdateUser Password Must Be Equal To New Password")

	//Testing DeleteUser
	dataManager.DeleteUser(userID, mockUsername, newPassword)
	user, err = dataManager.GetUserData(userID, config.DBConfigUserMetaDataCollection)
	assert.NoError(err, "GetUserData Returned Error")
	accountStatus := user[config.DBConfigAccountStatusField]
	assert.Equal(accountStatus, config.UserDataConfigAccountStatusDeleted, "Failed At DeleteUser Account Status Must Be Deleted")
}

func TestUserIDGen(t *testing.T) {
	assert := assert.New(t)
	uuid, err := generateUUID()
	assert.NoError(err, "Generate UUID Returned Error")
	assert.True(govalidator.IsUUIDv4(uuid))
}
