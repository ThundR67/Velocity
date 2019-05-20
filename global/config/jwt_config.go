package config

import "time"

//Here are all the config related to JWT

//Loading config manager
var jWTConfigManager = getConfigManager("jwt.config")

var (
	//JWTConfigSigningSecret This will be used to sign JWT
	JWTConfigSigningSecret = []byte("PrettySimpleAsOfNow")
	//JWTConfigUserIdentityField provides the user identity
	JWTConfigUserIdentityField = getStringConfig("jwt.userIdentityField", jWTConfigManager)
	//JWTConfigIsFreshField provides if the token is type fresh
	JWTConfigIsFreshField = getStringConfig("jwt.isFreshField", jWTConfigManager)
	//JWTConfigScopesField provides all the scopes
	JWTConfigScopesField = getStringConfig("jwt.scopesField", jWTConfigManager)
	//JWTConfigCreationUTCField is UTC of when token was created
	JWTConfigCreationUTCField = getStringConfig("jwt.creationUTCField", jWTConfigManager)
	//JWTConfigExpirationUTCField is UTC of when token will expire
	JWTConfigExpirationUTCField = getStringConfig("jwt.expirationUTCField", jWTConfigManager)
	//JWTConfigAccessTokenExpirationTimeMinutes is expiration time of accces token in minutes
	JWTConfigAccessTokenExpirationTimeMinutes = time.Duration(getIntConfig("jwt.accessTokenExpirationTimeMinutes", jWTConfigManager))
	//JWTConfigFreshAccessTokenExpirationTimeMinutes is expiration time of fresh accces token in minutes
	JWTConfigFreshAccessTokenExpirationTimeMinutes = time.Duration(getIntConfig("jwt.freshAccessTokenExpirationTimeMinutes", jWTConfigManager))
	//JWTConfigRefreshTokenExpirationTimeDays is expiration time of refresh token in days
	JWTConfigRefreshTokenExpirationTimeDays = time.Duration(getIntConfig("jwt.refreshTokenExpirationTimeDays", jWTConfigManager))
)
