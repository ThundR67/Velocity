package jwt

import (
	"testing"
	"time"

	"github.com/SonicRoshan/Velocity/global/config"
	goJwt "github.com/dgrijalva/jwt-go"
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
	accessToken, refreshToken := jwt.AccessAndRefreshTokens(id, scopes)

	testAccessAndRefresh(accessToken, refreshToken, id, scopes, t)

	accessToken, refreshToken, msg, err := jwt.RefreshTokens(refreshToken)
	testAccessAndRefresh(accessToken, refreshToken, id, scopes, t)
	assert.Zero(msg, "Message should be blank")
	assert.NoError(err, "RefreshTokens Returned Error")

	fresh := jwt.FreshToken(id)
	testToken(fresh,
		config.TokenTypeFresh,
		id,
		nil,
		config.JWTFreshAccessExpirationMinutes,
		t,
	)

	accessToken, refreshToken, msg, err = jwt.RefreshTokens("InvalidToken")
	assert.Error(err)
	assert.Zero(accessToken)
	assert.Zero(refreshToken)
	assert.Zero(msg)

	accessToken, _ = jwt.AccessAndRefreshTokens("test", nil)
	/*
		This should be invalid as we are passing accessToken
		to jwt.RefreshTokens which accepts refresh token not access token
	*/
	accessToken, refreshToken, msg, err = jwt.RefreshTokens(accessToken)
	assert.NoError(err)
	assert.Zero(accessToken)
	assert.Zero(refreshToken)
	assert.NotNil(msg)

	valid, claims, msg, err := jwt.ValidateToken("InvalidToken", config.TokenTypeAccess)
	assert.False(valid)
	assert.Zero(claims)
	assert.Zero(msg)
	assert.Error(err)

	accessToken, _ = jwt.AccessAndRefreshTokens("Test", nil)
	valid, claims, msg, err = jwt.ValidateToken(accessToken, "Invalid Type")
	assert.False(valid)
	assert.Zero(claims)
	assert.Zero(msg)
	assert.Error(err)
	assert.Equal(config.InvalidTokenMsg, err.Error())

	expiredAccessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWRlbnRpdHkiOiJUZXN0IiwiSXNGcmVzaCI6ZmFsc2UsIklzUmVmcmVzaCI6ZmFsc2UsIlNjb3BlcyI6bnVsbCwiQ3JlYXRpb25VVEMiOjE1Njg1NTgwMTAsIkV4cGlyYXRpb25VVEMiOjE1Njg1NTgwMTB9.cJJkt-N-Gl_8fEXDdke8dG96M43qK7F2RkMhhZrG8QQ"
	valid, claims, msg, err = jwt.ValidateToken(expiredAccessToken, config.TokenTypeAccess)
	assert.False(valid)
	assert.Zero(claims)
	assert.Equal(config.TokenExpiredMsg, msg)
	assert.NoError(err)

	invalidAccessToken := "eyJhbGciOiJQUzM4NCIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.MqF1AKsJkijKnfqEI3VA1OnzAL2S4eIpAuievMgD3tEFyFMU67gCbg-fxsc5dLrxNwdZEXs9h0kkicJZ70mp6p5vdv-j2ycDKBWg05Un4OhEl7lYcdIsCsB8QUPmstF-lQWnNqnq3wra1GynJrOXDL27qIaJnnQKlXuayFntBF0j-82jpuVdMaSXvk3OGaOM-7rCRsBcSPmocaAO-uWJEGPw_OWVaC5RRdWDroPi4YL4lTkDEC-KEvVkqCnFm_40C-T_siXquh5FVbpJjb3W2_YvcqfDRj44TsRrpVhk6ohsHMNeUad_cxnFnpolIKnaXq_COv35e9EgeQIPAbgIeg"
	_, err = goJwt.ParseWithClaims(invalidAccessToken, &claims, jwtKeyFunc)
	assert.Error(err)
}
