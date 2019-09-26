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
	userID, msg := dataManager.Add(mockUserData, mockUserExtraData)
	assert.Zero(msg)

	mockUserData2 := mockUserData
	mockUserData2.Username += "123"
	exists, msg := dataManager.doesUsernameOrEmailExists(mockUserData2)
	assert.Equal(msg, config.EmailExistMsg)
	assert.True(exists)

	mockUserData2 = mockUserData
	mockUserData2.Email = "123" + mockUserData2.Email
	exists, msg = dataManager.doesUsernameOrEmailExists(mockUserData2)
	assert.Equal(msg, config.UsernameExistMsg)
	assert.True(exists)

	//testing get
	var mainData config.UserMain
	dataManager.Get(userID, config.DBConfigUserMainDataCollection, &mainData)
	assert.NotZero(mainData)

	var extraData config.UserExtra
	dataManager.Get(userID, config.DBConfigUserExtraDataCollection, &extraData)
	assert.NotZero(extraData)

	var metaData config.UserMeta
	dataManager.Get(userID, config.DBConfigUserMetaDataCollection, &metaData)
	assert.NotZero(metaData)
	if !config.DebugMode {
		assert.Equal(
			config.UserDataConfigAccountStatusUnactivated,
			metaData.AccountStatus,
		)
	}

	//Testing activate
	msg = dataManager.Activate(mockUserData.Email)
	assert.Zero(msg)

	dataManager.Get(userID, config.DBConfigUserMetaDataCollection, &metaData)
	assert.Equal(
		config.UserDataConfigAccountStatusActive,
		metaData.AccountStatus,
	)

	//Testing auth user
	valid, _, _ := dataManager.Auth(mockUserData.Username, "", mockUserData.Password)
	assert.True(valid)
	valid, _, _ = dataManager.Auth(mockUserData.Username, "", "wrong-password")
	assert.False(valid)

	//Testing Update
	newPassword := "12345678"
	dataManager.Update(userID,
		config.UserMain{Password: newPassword},
		config.DBConfigUserMainDataCollection)

	//Testing Delete
	dataManager.Delete(userID, mockUserData.Username, newPassword)
	dataManager.Get(userID, config.DBConfigUserMetaDataCollection, &metaData)
	accountStatus := metaData.AccountStatus
	assert.Equal(
		config.UserDataConfigAccountStatusDeleted,
		accountStatus,
	)

	//Testing disconnect
	var newData config.UserMain
	dataManager.Disconnect()
	dataManager.Get(userID, config.DBConfigUserMainDataCollection, &newData)
	assert.Zero(newData)
}

func TestUserIDGen(t *testing.T) {
	assert := assert.New(t)
	uuid, err := generateUUID()
	assert.NoError(err)
	assert.True(govalidator.IsUUIDv4(uuid))
}
