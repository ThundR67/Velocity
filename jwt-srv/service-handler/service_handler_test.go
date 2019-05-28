package handler

import (
	"testing"

	"github.com/SonicRoshan/Velocity/global/clients"
	"github.com/SonicRoshan/Velocity/global/config"
	micro "github.com/micro/go-micro"
	"github.com/stretchr/testify/assert"
)

func testToken(
	token,
	tokenType,
	id string,
	scopes []string,
	assert *assert.Assertions,
	client clients.JWTClient) {

	valid, idInToken, scopesInToken, err := client.ValidateToken(token, tokenType)
	assert.True(valid, "Token Should Be Valid")
	assert.NoError(err, "ValidateToken Returned Error")
	assert.Equal(id, idInToken, "IDs Should Match")
	assert.Equal(scopes, scopesInToken, "Scopes Should Match")
}

func TestJWTServiceHandler(t *testing.T) {
	assert := assert.New(t)
	service := micro.NewService(micro.Name("TestService"))
	jwtClient := clients.NewJWTClient(service)

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
