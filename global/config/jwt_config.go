package config

import "time"

//Here are all the config related to JWT

//Loading config manager
var jWTConfigManager = getConfigManager("jwt.config")

var (
	//JWTSecret is used sign jwts
	JWTSecret = []byte("Cause I am not listening")

	//JWTAccessExpirationMinutes is used to add expiration time to access token
	JWTAccessExpirationMinutes = time.Minute * time.Duration(getIntConfig("jwt.accessTokenExpirationTimeMinutes", jWTConfigManager))
	//JWTFreshAccessExpirationMinutes is used to add expiration time to fresh token
	JWTFreshAccessExpirationMinutes = time.Minute * time.Duration(getIntConfig("jwt.freshAccessTokenExpirationTimeMinutes", jWTConfigManager))
	//JWTRefreshExpirationDays is used to add expiration time to refresh access token
	JWTRefreshExpirationDays = time.Hour * 24 * time.Duration(getIntConfig("jwt.refreshTokenExpirationTimeDays", jWTConfigManager))
)

const (
	//TokenTypeAccess is used to specify access token while validating token
	TokenTypeAccess = "access"
	//TokenTypeRefresh is used to specify refresh token while validating token
	TokenTypeRefresh = "refresh"
	//TokenTypeFresh is used to specify fresh token while validating token
	TokenTypeFresh = "fresh"
)
