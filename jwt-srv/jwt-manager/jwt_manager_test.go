package jwtmanager

import (
	"testing"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/stretchr/testify/assert"
)

func match(a []interface{}, b []string) bool {
	l := len(a)
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < l; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

//TestJWTManager test jwt-manager.go
func TestJWTManager(t *testing.T) {
	assert := assert.New(t)
	manager := JWTManager{}
	scopes := []string{"read", "write", "delete"}

	//Testing generation of access and refresh token
	accessToken, refreshToken, _ := manager.GenerateAccessAndRefreshToken("test", scopes)
	//testing ValidateToken
	valid, claims, _, err := manager.ValidateToken(accessToken)

	assert.True(match(claims[config.JWTConfigScopesField].([]interface{}), scopes), "Scopes in claims should match scopes entered")
	assert.NoError(err, "Error While Validating Token")
	assert.True(valid, "Access Token Should Be Valid")

	valid, claims, _, err = manager.ValidateToken(refreshToken)

	assert.True(match(claims[config.JWTConfigScopesField].([]interface{}), scopes), "Scopes in claims should match scopes entered")
	assert.NoError(err, "Error While Validating Token")
	assert.True(valid, "Refresh Token Should Be Valid")

	valid, _, _, _ = manager.ValidateToken("randomstuff")
	assert.False(valid, "Random Token Should Not Be Valid")

	//testing generation of refresh token and access token thru refresh token
	accessToken, refreshToken, _, err = manager.GenerateAccessAndRefreshTokenBasedOnRefreshToken(refreshToken)
	accessValid, _, _, err := manager.ValidateToken(accessToken)
	assert.NoError(err, "Error While Validating Token")
	refreshValid, _, _, err := manager.ValidateToken(refreshToken)
	assert.NoError(err, "Error While Validating Token")
	assert.True(accessValid && refreshValid, "New Access And Refresh Token Should Be Valid")

	//testing generation of fresh token
	freshAccessToken, err := manager.GenerateFreshAccesToken("test")
	assert.NoError(err, "Error While Generating Fresh Acces Token")

	//testing validation of fresh access token
	freshValid, _, err := manager.ValidateFreshAccessToken(freshAccessToken)
	assert.NoError(err, "Error While Validating Fresh Access Token")
	assert.True(freshValid, "Fresh Access Token Should Be Valid")
	randomValid, _, _ := manager.ValidateFreshAccessToken("random")
	assert.False(randomValid, "Random Data Should Not Be Valid")
}
