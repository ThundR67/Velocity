package users

import (
	"encoding/json"
	"fmt"
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

func validateMain(key, value string) bool {
	mockUsername := strings.ToLower(generateRandomString(10))
	mockPassword := generateRandomString(10)
	toValidateStr := fmt.Sprintf(`{"%s":"%s"}`, key, value)
	toValidate := config.UserMain{
		Username: mockUsername,
		Password: mockPassword,
		Email:    mockUsername + "@gmail.com",
	}
	err := json.Unmarshal([]byte(toValidateStr), &toValidate)
	if err != nil {
		panic(fmt.Sprintf("Key %s Val %s Err %s", key, value, err.Error()))
	}
	return validateUserMainData(toValidate, false)
}

func validateExtra(key, value string) bool {
	mockUsername := strings.ToLower(generateRandomString(10))
	toValidateStr := fmt.Sprintf(`{"%s":"%s"}`, key, value)
	toValidate := config.UserExtra{
		Gender:      "male",
		FirstName:   mockUsername,
		LastName:    mockUsername,
		BirthdayUTC: int64(864466669), //A Timestamp of year 1997
	}
	err := json.Unmarshal([]byte(toValidateStr), &toValidate)
	if err != nil {
		panic(fmt.Sprintf("Key %s Val %s Err %s", key, value, err.Error()))
	}
	return validateUserExtraData(toValidate, false)
}

func TestUtils(t *testing.T) {
	assert := assert.New(t)

	metaData := generateUserMetaData()
	assert.Equal(config.UserDataConfigAccountStatusUnactivated,
		metaData.AccountStatus,
		"GenerateMetaData Should Return Data With Account Status Unactivated")
	assert.WithinDuration(time.Unix(metaData.AccountCreationUTC, 0),
		time.Now(),
		time.Second*5,
		"Account Creation Time Is Incorrect")
}

//TestLowLevelUsers is used to test Users struct
func TestLowLevelUsers(t *testing.T) {
	assert := assert.New(t)

	dataManager := Users{DBName: "TestDB"}
	err := dataManager.Init()
	assert.NoError(err, "Initializing Caused Error")

	mockUsername := strings.ToLower(generateRandomString(10))
	mockPassword := generateRandomString(10)

	mockUserData := config.UserMain{
		Username: mockUsername,
		Password: mockPassword,
		Email:    mockUsername + "@gmail.com",
	}

	mockUserExtraData := config.UserExtra{
		Gender:      "male",
		FirstName:   mockUsername,
		LastName:    mockUsername,
		BirthdayUTC: int64(864466669), //A Timestamp of year 1997
	}

	//Testing add
	userID, msg, err := dataManager.Add(mockUserData, mockUserExtraData)
	assert.NoError(err, "Add Returned Error")
	assert.Equal(msg, "", "Message Should Be Blank")

	//testing get
	var mainData config.UserMain
	err = dataManager.Get(userID, config.DBConfigUserMainDataCollection, &mainData)
	assert.NoError(err, "Get Returned Error For Getting Main Data")
	assert.True(mainData != (config.UserMain{}))

	var extraData config.UserExtra
	err = dataManager.Get(userID, config.DBConfigUserExtraDataCollection, &extraData)
	assert.NoError(err, "Get Returned Error For Getting Extra Data")
	assert.True(extraData != (config.UserExtra{}))

	var metaData config.UserMeta
	err = dataManager.Get(userID, config.DBConfigUserMetaDataCollection, &metaData)
	assert.NoError(err, "Get Returned Error For Getting Meta Data")
	assert.True(metaData != (config.UserMeta{}))
	if !config.DebugMode {
		assert.Equal(
			config.UserDataConfigAccountStatusUnactivated,
			metaData.AccountStatus,
			"Account Status Should Be Unactive",
		)
	}

	//Testing activate
	msg, err = dataManager.Activate(mockUsername + "@gmail.com")
	assert.NoError(err, "Add Returned Error")
	assert.Equal(msg, "", "Message Should Be Blank")

	dataManager.Get(userID, config.DBConfigUserMetaDataCollection, &metaData)
	assert.Equal(
		config.UserDataConfigAccountStatusActive,
		metaData.AccountStatus,
		"Account Status Should Be Activate",
	)

	//testing get with field names
	err = dataManager.Get(userID, config.DBConfigUserMetaDataCollection, &metaData)
	assert.NoError(err, "Get Returned Error")

	//Testing auth user
	valid, _, _, err := dataManager.Auth(mockUsername, "", mockPassword)
	assert.NoError(err, "Auth Returned Error")
	assert.True(valid, "Failed At Auth For Checking With Correct Password")
	valid, _, _, _ = dataManager.Auth(mockUsername, "", "wrong-password")
	assert.False(valid, "Failed At Auth For Checking With Incorrect Password")

	//Testing Update
	newPassword := generateRandomString(10)
	err = dataManager.Update(userID,
		config.UserMain{Password: newPassword},
		config.DBConfigUserMainDataCollection)

	assert.NoError(err, "UpdateData Returned Error")

	//Testing Delete
	_, err = dataManager.Delete(userID, mockUsername, newPassword)
	assert.NoError(err, "Delete Returned Error")
	err = dataManager.Get(userID, config.DBConfigUserMetaDataCollection, &metaData)
	assert.NoError(err, "Get Returned Error")
	accountStatus := metaData.AccountStatus
	assert.Equal(config.UserDataConfigAccountStatusDeleted,
		accountStatus,
		"Failed At Delete Account Status Must Be Deleted")
}

func TestUserIDGen(t *testing.T) {
	assert := assert.New(t)
	uuid, err := generateUUID()
	assert.NoError(err, "Generate UUID Returned Error")
	assert.True(govalidator.IsUUIDv4(uuid))
}

func TestValidators(t *testing.T) {
	assert := assert.New(t)

	mainData := map[string][2]string{
		"Username": [2]string{"testing", "ta"}, //where first is valid and second is not
		"Password": [2]string{"testingPassword", "ta"},
		"Email":    [2]string{"sonicroshan122@gmail.com", "thisemailisnotvalid"},
	}

	extraData := map[string][2]string{
		"FirstName": [2]string{"testing", "ta"},
		"LastName":  [2]string{"testing", "ta"},
		"Gender":    [2]string{"male", "notvalidgender"},
	}

	var valid, invalid string
	for key, value := range mainData {
		valid = value[0]
		invalid = value[1]
		assert.Truef(validateMain(key, valid), "%s Key With Val %s Should Be Valid", key, valid)
		assert.Falsef(validateMain(key, invalid), "%s Key With Val %s Should Not Be Valid", key, invalid)
	}

	for key, value := range extraData {
		valid = value[0]
		invalid = value[1]
		assert.Truef(validateExtra(key, valid), "%s Key With Val %s Should Be Valid", key, valid)
		assert.Falsef(validateExtra(key, invalid), "%s Key With Val %s Should Not Be Valid", key, invalid)
	}
	//Testing bday
	mockUsername := strings.ToLower(generateRandomString(10))
	validBirthday := config.UserExtra{
		Gender:      "male",
		FirstName:   mockUsername,
		LastName:    mockUsername,
		BirthdayUTC: int64(864466669), //A Timestamp of year 1997
	}
	invalidBirthday := config.UserExtra{
		Gender:      "male",
		FirstName:   mockUsername,
		LastName:    mockUsername,
		BirthdayUTC: time.Now().Unix(), //A Timestamp of year 1997
	}

	assert.True(validateUserExtraData(validBirthday, false))
	assert.False(validateUserExtraData(invalidBirthday, false))

}
