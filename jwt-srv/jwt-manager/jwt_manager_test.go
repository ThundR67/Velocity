package jwtmanager

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//TestJWTManager test jwt-manager.go
func TestJWTManager(t *testing.T) {
	manager := JWTManager{}

	//Testing generation of access and refresh token
	accessToken, refreshToken, _ := manager.GenerateAccessAndRefreshToken("test", []string{"read", "write", "delete"})
	//testing ValidateToken
	valid, _, _ := manager.ValidateToken(accessToken)
	assert.True(t, valid, "Access Token Should Be Valid")
	valid, _, _ = manager.ValidateToken(refreshToken)
	assert.True(t, valid, "Refresh Token Should Be Valid")
	valid, _, _ = manager.ValidateToken("randomstuff")
	assert.False(t, valid, "Random Token Should Not Be Valid")

	//testing generation of refresh token and access token thru refresh token
	accessToken, refreshToken, _ = manager.GenerateAccessAndRefreshTokenBasedOnRefreshToken(refreshToken)
	accessValid, _, _ := manager.ValidateToken(accessToken)
	refreshValid, _, _ := manager.ValidateToken(refreshToken)
	assert.True(t, accessValid && refreshValid, "New Access And Refresh Token Should Be Valid")

	//testing generation of fresh token
	freshAccessToken, _ := manager.GenerateFreshAccesToken("test")

	//testing validation of fresh access token
	freshValid, _ := manager.ValidateFreshAccessToken(freshAccessToken)
	fmt.Println(freshAccessToken)
	assert.True(t, freshValid, "Fresh Access Token Should Be Valid")
	randomValid, _ := manager.ValidateFreshAccessToken("random")
	assert.False(t, randomValid, "Random Data Should Not Be Valid")
}
