package clients

import (
	"math/rand"
	"strings"
	"testing"
	"time"

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

func testToken(
	token,
	tokenType,
	id string,
	scopes []string,
	assert *assert.Assertions,
	client JWTClient) {

	valid, idInToken, scopesInToken, err := client.ValidateToken(token, tokenType)
	assert.True(valid, "Token Should Be Valid")
	assert.NoError(err, "ValidateToken Returned Error")
	assert.Equal(id, idInToken, "IDs Should Match")
	assert.Equal(scopes, scopesInToken, "Scopes Should Match")
}



func TestJWTClient(t *testing.T) {
	assert := assert.New(t)
	service := micro.NewService(micro.Name("TestService"))
	jwtClient := NewJWTClient(service)

	id := "TestID"
	scopes := []string{"READ", "WRITE", "DELETE"}

	access, refresh, msg, err := jwtClient.AccessAndRefreshTokens(id, scopes)
	assert.Zero(msg, "Message Should Be Blank")
	assert.NoError(err, "AccessAndRefreshTokens Returned Error")

	testToken(access, config.TokenTypeAccess, id, scopes, assert, jwtClient)
	testToken(refresh, config.TokenTypeRefresh, id, scopes, assert, jwtClient)

	access, refresh, msg, err = jwtClient.RefreshTokens(refresh)
	assert.Zero(msg, "Message Should Be Blank")
	assert.NoError(err, "RefreshTokens Returned Error")

	testToken(access, config.TokenTypeAccess, id, scopes, assert, jwtClient)
	testToken(refresh, config.TokenTypeRefresh, id, scopes, assert, jwtClient)

	fresh, err := jwtClient.FreshToken(id)
	assert.NoError(err, "FreshToken Returned Error")

	testToken(fresh, config.TokenTypeFresh, id, nil, assert, jwtClient)
}


func TestUsersClient(t *testing.T) {
	assert := assert.New(t)
	service := micro.NewService(micro.Name("test"))
	client := NewUsersClient(service)

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

	//Testing Add
	userID, msg, err := client.Add(mockUserData, mockUserExtraData)
	assert.NoError(err, "Adding User Returned Error")
	assert.Equal("", msg, "Message By Add Should Be Blank")

	//Testing Get
	data, msg, err := client.Get(userID)
	assert.NoError(err, "Getting User Returned Error")
	assert.Equal("", msg, "Message By Get Should Be Blank")
	assert.Equal(mockUsername, data.Username, "Username Should Match")

	//Testing Update
	updatedEmail := generateRandomString(10) + "@gmail.com"
	update := config.UserMain{
		Email: updatedEmail,
	}
	err = client.Update(userID, update)
	assert.NoError(err, "UpdateUser Returned Error")
	data, _, _ = client.Get(userID)
	assert.Equal(updatedEmail, data.Email, "Updated Email Should Match Email In DB")

	//Testing Auth
	id, msg, err := client.Auth(mockUsername, mockPassword)
	assert.NoError(err, "Auth User Returned Error")
	assert.Equal("", msg, "Message By Auth Should Be Blank")
	assert.Equal(userID, id, "Ids should match")

}
