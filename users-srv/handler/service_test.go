package handler

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/SonicRoshan/Velocity/global/clients"
	"github.com/SonicRoshan/Velocity/global/config"
	micro "github.com/micro/go-micro"
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

func TestUsersService(t *testing.T) {
	assert := assert.New(t)
	service := micro.NewService(micro.Name("test"))
	client := clients.NewUsersClient(service)

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

	//Testing Add
	userID, msg, err := client.Add(mockUserData, mockUserExtraData)
	assert.NoError(err, "Adding User Returned Error")
	assert.Equal("", msg, "Message By Add Should Be Blank")

	//Testing Get
	data, msg, err := client.Get(userID, config.DBConfigUserMainDataCollection)
	assert.NoError(err, "Getting User Returned Error")
	assert.Equal("", msg, "Message By Get Should Be Blank")
	assert.Equal(mockUsername, data.Username, "Username Should Match")

	//Testing Auth
	id, msg, err := client.Auth(mockUsername, mockPassword)
	assert.NoError(err, "Auth User Returned Error")
	assert.Equal("", msg, "Message By Auth Should Be Blank")
	assert.Equal(userID, id, "Ids should match")

}
