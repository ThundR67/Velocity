package userdatamanager

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/SonicRoshan/Velocity/user-data-srv/config"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Creating A MongoDB Client To TestDB Where All Testing Will Happend
var client, _ = mongo.NewClient(options.Client().ApplyURI(config.ConfigMongoDBAddress))
var _ = func() string { fmt.Println(config.ConfigMongoDBAddress); return "h" }()
var err = client.Connect(context.TODO())
var collection = client.Database("TestDB").
	Collection(config.ConfigUserDataCollection)

//Creating dataManager With collection Created To TestDB
var dataManager = UserDataManager{Collection: collection, Ctx: context.TODO()}

func generateRandomString(length int) string {
	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = config.ConfigUserIDCharset[seededRand.Intn(len(config.ConfigUserIDCharset))]
	}
	return string(b)
}

//TestLowLevelUserDataManager Will Test Low Level user data manager
func TestLowLevelUserDataManager(t *testing.T) {
	//Creating mock Data For Testing
	mockUsername := generateRandomString(10)
	mockPassword := generateRandomString(10)

	mockUserData := map[string]interface{}{config.ConfigUsernameField: mockUsername,
		config.ConfigPasswordField:      mockPassword,
		config.ConfigEmailField:         mockUsername,
		config.ConfigUserExtraDataField: map[string]interface{}{}}

	//Testing AddUser
	userID, err := dataManager.AddUser(mockUserData)
	if err != nil {
		t.Log(err)
		t.Error("Failed At AddUser")
	}

	//Testing GetUser
	user, _ := dataManager.GetUser(userID)
	assert.Equal(t, mockUsername, user[config.ConfigUsernameField], "Failed At GetUser, Mock Username And Username In DB Must Be Equal")

	//Testing GetUserByUsernameOrEmail
	user, _ = dataManager.GetUserByUsernameOrEmail(mockUsername, "")
	assert.Equal(t, mockUsername, user[config.ConfigUsernameField], "Failed At GetUserByUsernameOrEmail, Mock Username And Username In DB Must Be Equal")

	//Testing auth user
	valid, _, err := dataManager.AuthUser(mockUsername, "", mockPassword)
	assert.True(t, valid, "Failed At AuthUser For Checking With Correct Password")
	valid, _, _ = dataManager.AuthUser(mockUsername, "", "wrong-password")
	assert.True(t, !valid, "Failed At AuthUser For Checking With Incorrect Password")

	//Testing UpdateUser
	newPassword := generateRandomString(10)
	dataManager.UpdateUser(userID, config.ConfigPasswordField, newPassword)
	user, _ = dataManager.GetUserByUsernameOrEmail(mockUsername, "", true)
	assert.Equal(t, newPassword, user[config.ConfigPasswordField], "Failed At UpdateUser Password Must Be Equal To New Password")

	//Testing DeleteUser
	dataManager.DeleteUser(userID, mockUsername, newPassword)
	user, _ = dataManager.GetUser(userID)
	accountStatus := user[config.ConfigUserExtraDataField].(primitive.M)[config.ConfigAccountStatusField]
	assert.Equal(t, accountStatus, config.ConfigAccountStatusDeleted, "Failed At DeleteUser Account Status Must Be Deleted")

}
