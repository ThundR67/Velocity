package config

import (
	"os"
	"strings"
	"time"

	goup "github.com/ufoscout/go-up"
)

func getConfigFilePath() string {
	dir, _ := os.Getwd()
	split := strings.Split(dir, "\\")
	if split[len(split)-1] == "jwt-manager" {
		return "../config/main.config"
	}
	return "config/main.config"
}

//Loading Up Config File
var configurations, _ = goup.NewGoUp().
	AddFile(getConfigFilePath(), false).
	Build()

var (
	//ConfigSigningSecret This will be used to sign JWT
	ConfigSigningSecret = []byte("PrettySimpleAsOfNow")
	//ConfigUserIdentityField provides the user identity
	ConfigUserIdentityField, _ = configurations.GetStringOrFail("jwt.userIdentityField")
	//ConfigIsFreshField provides if the token is type fresh
	ConfigIsFreshField = configurations.GetString("jwt.isFreshField")
	//ConfigScopesField provides all the scopes
	ConfigScopesField = configurations.GetString("jwt.scopesField")
	//ConfigCreationUTCField is UTC of when token was created
	ConfigCreationUTCField = configurations.GetString("jwt.creationUTCField")
	//ConfigExpirationUTCField is UTC of when token will expire
	ConfigExpirationUTCField = configurations.GetString("jwt.expirationUTCField")
	//ConfigAccessTokenExpirationTimeMinutes is expiration time of accces token in minutes
	ConfigAccessTokenExpirationTimeMinutes = time.Duration(configurations.GetInt("jwt.accessTokenExpirationTimeMinutes"))
	//ConfigFreshAccessTokenExpirationTimeMinutes is expiration time of fresh accces token in minutes
	ConfigFreshAccessTokenExpirationTimeMinutes = time.Duration(configurations.GetInt("jwt.freshAccessTokenExpirationTimeMinutes"))
	//ConfigRefreshTokenExpirationTimeDays is expiration time of refresh token in days
	ConfigRefreshTokenExpirationTimeDays = time.Duration(configurations.GetInt("jwt.refreshTokenExpirationTimeDays"))
)
