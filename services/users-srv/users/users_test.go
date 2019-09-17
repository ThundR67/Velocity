package users

import (
	"testing"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/SonicRoshan/Velocity/global/utils"
	"github.com/asaskevich/govalidator"
	"github.com/stretchr/testify/assert"
)

//TestLowLevelUsers is used to test Users struct
func TestLowLevelUsers(t *testing.T) {
	assert := assert.New(t)

	dataManager := Users{DBName: "TestDB"}
	dataManager.Init()

	mockUserData, mockUserExtraData := utils.GetMockUserData()

	//Testing add
	userID, msg, err := dataManager.Add(mockUserData, mockUserExtraData)
	assert.NoError(err)
	assert.Zero(msg)

	mockUserData2 := mockUserData
	mockUserData2.Username += "123"
	exists, msg, err := dataManager.doesUsernameOrEmailExists(mockUserData2)
	assert.NoError(err)
	assert.Equal(msg, config.EmailExistMsg)
	assert.True(exists)

	mockUserData2 = mockUserData
	mockUserData2.Email = "123" + mockUserData2.Email
	exists, msg, err = dataManager.doesUsernameOrEmailExists(mockUserData2)
	assert.NoError(err)
	assert.Equal(msg, config.UsernameExistMsg)
	assert.True(exists)

	//testing get
	var mainData config.UserMain
	err = dataManager.Get(userID, config.DBConfigUserMainDataCollection, &mainData)
	assert.NoError(err)
	assert.NotZero(mainData)

	var extraData config.UserExtra
	err = dataManager.Get(userID, config.DBConfigUserExtraDataCollection, &extraData)
	assert.NoError(err)
	assert.NotZero(extraData)

	var metaData config.UserMeta
	err = dataManager.Get(userID, config.DBConfigUserMetaDataCollection, &metaData)
	assert.NoError(err)
	assert.NotZero(metaData)
	if !config.DebugMode {
		assert.Equal(
			config.UserDataConfigAccountStatusUnactivated,
			metaData.AccountStatus,
		)
	}

	//Testing activate
	msg, err = dataManager.Activate(mockUserData.Email)
	assert.NoError(err)
	assert.Zero(msg)

	dataManager.Get(userID, config.DBConfigUserMetaDataCollection, &metaData)
	assert.Equal(
		config.UserDataConfigAccountStatusActive,
		metaData.AccountStatus,
	)

	//testing get with field names
	err = dataManager.Get(userID, config.DBConfigUserMetaDataCollection, &metaData)
	assert.NoError(err, "Get Returned Error")

	//Testing auth user
	valid, _, _, err := dataManager.Auth(mockUserData.Username, "", mockUserData.Password)
	assert.NoError(err)
	assert.True(valid)
	valid, _, _, _ = dataManager.Auth(mockUserData.Username, "", "wrong-password")
	assert.False(valid)

	//Testing Update
	newPassword := "12345678"
	err = dataManager.Update(userID,
		config.UserMain{Password: newPassword},
		config.DBConfigUserMainDataCollection)

	assert.NoError(err)

	//Testing Delete
	_, err = dataManager.Delete(userID, mockUserData.Username, newPassword)
	assert.NoError(err)
	err = dataManager.Get(userID, config.DBConfigUserMetaDataCollection, &metaData)
	assert.NoError(err)
	accountStatus := metaData.AccountStatus
	assert.Equal(
		config.UserDataConfigAccountStatusDeleted,
		accountStatus,
	)

	//Testing disconnect
	var newData config.UserMain
	dataManager.Disconnect()
	err = dataManager.Get(userID, config.DBConfigUserMainDataCollection, &newData)
	assert.Error(err)
	assert.Zero(newData)
}

func TestUserIDGen(t *testing.T) {
	assert := assert.New(t)
	uuid, err := generateUUID()
	assert.NoError(err)
	assert.True(govalidator.IsUUIDv4(uuid))
}
