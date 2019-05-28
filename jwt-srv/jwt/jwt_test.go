package jwt

import (
	"testing"
	"time"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/stretchr/testify/assert"
)

func testToken(
	token,
	tokenType,
	id string,
	scopes []string,
	expirationTime time.Duration,
	t *testing.T) {

	assert := assert.New(t)
	jwt := JWT{}

	valid, claims, msg, err := jwt.ValidateToken(token, tokenType)
	assert.NoError(err, "Validate Token Returned Error")
	assert.Zero(msg, "Message should be blank")
	assert.NotNil(claims)
	assert.True(valid, "Token Should Be Valid")

	assert.WithinDuration(
		time.Now(),
		time.Unix(claims.CreationUTC, 0),
		time.Second,
		"Creation Time Is Incorrect")

	assert.WithinDuration(
		time.Now().Add(expirationTime),
		time.Unix(claims.ExpirationUTC, 0),
		time.Second,
		"Expiration Time Is Incorrect")

	assert.Equal(scopes, claims.Scopes, "Scopes Should Match")
}

func testAccessAndRefresh(access, refresh, id string, scopes []string, t *testing.T) {
	testToken(access,
		config.TokenTypeAccess,
		id,
		scopes,
		config.JWTAccessExpirationMinutes,
		t)

	testToken(refresh,
		config.TokenTypeRefresh,
		id,
		scopes,
		config.JWTRefreshExpirationDays,
		t)
}

func TestClaimsGen(t *testing.T) {
	assert := assert.New(t)

	scopes := []string{"read", "write", "delete"}
	id := "testing"

	accessClaims := accessTokenClaims(id, scopes)
	assert.True(!accessClaims.IsFresh && !accessClaims.IsRefresh,
		"Access token claims should neither be isFresh Nor IsRefresh")

	refreshClaims := refreshTokenClaims(id, scopes)
	assert.True(refreshClaims.IsRefresh && !refreshClaims.IsFresh,
		"Refresh Claims Should Be IsRefresh")

	freshClaims := freshTokenClaims(id)
	assert.True(freshClaims.IsFresh && !freshClaims.IsRefresh,
		"Fresh Claims Should Be IsRefresh")
}

func TestJWT(t *testing.T) {
	assert := assert.New(t)
	jwt := JWT{}
	scopes := []string{"read", "write", "delete"}
	id := "testing"

	//Testing access and refresh token generation
	accessToken, refreshToken, err := jwt.AccessAndRefreshTokens(id, scopes)
	assert.NoError(err, "Generating Access And Refresh Token Returned Error")

	testAccessAndRefresh(accessToken, refreshToken, id, scopes, t)

	accessToken, refreshToken, msg, err := jwt.RefreshTokens(refreshToken)
	testAccessAndRefresh(accessToken, refreshToken, id, scopes, t)
	assert.Zero(msg, "Message should be blank")
	assert.NoError(err, "RefreshTokens Returned Error")

	fresh, err := jwt.FreshToken(id)
	assert.NoError(err, "Fresh Token Returned Error")
	testToken(fresh,
		config.TokenTypeFresh,
		id,
		nil,
		config.JWTFreshAccessExpirationMinutes,
		t)

}
