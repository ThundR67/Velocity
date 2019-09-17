package clients

import (
	"testing"

	"github.com/SonicRoshan/Velocity/global/config"
	micro "github.com/micro/go-micro"
	"github.com/stretchr/testify/assert"
)

func TestJWTClient(t *testing.T) {
	assert := assert.New(t)
	service := micro.NewService(micro.Name("TestService"))
	jwtClient := NewJWTClient(service)

	id := "TestID"
	scopes := []string{"READ", "WRITE", "DELETE"}

	access, refresh, msg, err := jwtClient.AccessAndRefreshTokens(id, scopes)
	assert.Zero(msg, "Message Should Be Blank")
	assert.NotZero(access)
	assert.NotZero(refresh)
	assert.NoError(err, "AccessAndRefreshTokens Returned Error")

	valid, id, scopes, err := jwtClient.ValidateToken(access, config.TokenTypeAccess)
	assert.True(valid)
	assert.NotZero(id)
	assert.NotZero(id)
	assert.NoError(err)

	access, refresh, msg, err = jwtClient.RefreshTokens(refresh)
	assert.Zero(msg)
	assert.NotZero(access)
	assert.NotZero(refresh)
	assert.NoError(err)

	access, refresh, msg, err = jwtClient.RefreshTokens("InvalidToken")
	assert.Zero(msg)
	assert.Error(err)

	fresh := jwtClient.FreshToken(id)
	assert.NotZero(fresh)

	valid, id, scopes, err = jwtClient.ValidateToken("SomeToken", "InvalidToken")
	assert.False(valid)
	assert.Zero(id)
	assert.Zero(scopes)
	assert.Error(err)
}
