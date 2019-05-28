package users

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

//TestLowLevelUsers is used to test Users struct
func TestLowLevelUsers(t *testing.T) {
	assert := assert.New(t)

	dataManager := Users{DBName: "TestDB"}
	err := dataManager.Init()
	assert.NoError(err, "Initializing Caused Error")

	mockUsername := strings.ToLower(generateRandomString(10))
	mockPassword := generateRandomString(10)

	mockUserData := config.UserType{
		Username: mockUsername,
		Password: mockPassword,
		Email:    mockUsername + "@gmail.com",
	}

	mockUserExtraData := config.UserType{
		Gender:      "male",
		FirstName:   mockUsername,
		LastName:    mockUsername,
		BirthdayUTC: int64(864466669), //A Timestamp of year 1997
	}

	//Testing add
	userID, msg, err := dataManager.Add(mockUserData, mockUserExtraData)
	assert.NoError(err, "Add Returned Error")
	assert.Equal(msg, "", "Message Should Be Blank")

	//testing of all three collection got the data
	data, err := dataManager.Get(userID, config.DBConfigUserMainDataCollection)
	assert.NoError(err, "GetData Returned Error")
	assert.True(data != (config.UserType{}))

	data, err = dataManager.Get(userID, config.DBConfigUserExtraDataCollection)
	assert.NoError(err, "GetData Returned Error")
	assert.True(data != (config.UserType{}))

	data, err = dataManager.Get(userID, config.DBConfigUserMetaDataCollection)
	assert.NoError(err, "GetData Returned Error")
	assert.True(data != (config.UserType{}))

	//Testing auth user
	valid, _, _, err := dataManager.Auth(mockUsername, "", mockPassword)
	assert.NoError(err, "Auth Returned Error")
	assert.True(valid, "Failed At Auth For Checking With Correct Password")
	valid, _, _, _ = dataManager.Auth(mockUsername, "", "wrong-password")
	assert.False(valid, "Failed At Auth For Checking With Incorrect Password")

	//Testing Update
	newPassword := generateRandomString(10)

	err = dataManager.Update(userID, config.UserType{Password: newPassword}, config.DBConfigUserMainDataCollection)
	assert.NoError(err, "UpdateData Returned Error")

	data, err = dataManager.Get(userID, config.DBConfigUserMainDataCollection)
	assert.NoError(err, "Get Returned Error")
	assert.Equal(newPassword, data.Password, "Failed At Update Password Must Be Equal To New Password")

	//Testing Delete
	_, err = dataManager.Delete(userID, mockUsername, newPassword)
	assert.NoError(err, "Delete Returned Error")
	user, err := dataManager.Get(userID, config.DBConfigUserMetaDataCollection)
	assert.NoError(err, "Get Returned Error")
	accountStatus := user.AccountStatus
	assert.Equal(config.UserDataConfigAccountStatusDeleted, accountStatus, "Failed At Delete Account Status Must Be Deleted")
}

func TestUserIDGen(t *testing.T) {
	assert := assert.New(t)
	uuid, err := generateUUID()
	assert.NoError(err, "Generate UUID Returned Error")
	assert.True(govalidator.IsUUIDv4(uuid))
}
