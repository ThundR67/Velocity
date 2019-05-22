package jwtmanager

import (
	"fmt"
	"time"

	"github.com/SonicRoshan/Velocity/global/config"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

//timestampFromFloat64 converts float64 to time.time
func timeFromFloat64(ts float64) time.Time {
	secs := int64(ts)
	nsecs := int64((ts - float64(secs)) * 1e9)
	return time.Unix(secs, nsecs)
}

//SliceInterfaceToString []interface{} to []string
func SliceInterfaceToString(slice []interface{}) []string {
	output := []string{}
	for _, val := range slice {
		output = append(output, val.(string))
	}
	return output
}

//JWTManager is a low level jason web token manager
type JWTManager struct{}

//isExpired checks if a time.Duration is expired
func (jwtManager JWTManager) isExpired(expirationTime time.Time) bool {
	return time.Now().After(expirationTime)
}

//generateClaimsForFreshToken generates claims for fresh acces token
func (jwtManager JWTManager) generateClaimsForFreshToken(userIdentity string) jwt.MapClaims {
	currentTime := time.Now()
	freshAccessTokenExpirationTime := currentTime.Add(time.Minute * config.JWTConfigFreshAccessTokenExpirationTimeMinutes)
	return jwt.MapClaims{
		config.JWTConfigUserIdentityField:  userIdentity,
		config.JWTConfigIsFreshField:       true,
		config.JWTConfigCreationUTCField:   currentTime.Unix(),
		config.JWTConfigExpirationUTCField: freshAccessTokenExpirationTime.Unix(),
	}
}

//generateClaims for acces and refresh token
func (jwtManager JWTManager) generateClaims(userIdentity string, scopes []string) (jwt.MapClaims, jwt.MapClaims) {
	currentTime := time.Now()
	accessTokenExpirationTime := currentTime.Add(time.Minute * config.JWTConfigAccessTokenExpirationTimeMinutes)
	refreshTokenExpirationTime := currentTime.Add(time.Hour * 24 * config.JWTConfigRefreshTokenExpirationTimeDays)
	accessTokenClaims := jwt.MapClaims{
		config.JWTConfigUserIdentityField:  userIdentity,
		config.JWTConfigIsFreshField:       false,
		config.JWTConfigScopesField:        scopes,
		config.JWTConfigCreationUTCField:   currentTime.Unix(),
		config.JWTConfigExpirationUTCField: accessTokenExpirationTime.Unix(),
	}

	refreshTokenClaims := jwt.MapClaims{
		config.JWTConfigUserIdentityField:  userIdentity,
		config.JWTConfigScopesField:        scopes,
		config.JWTConfigCreationUTCField:   currentTime,
		config.JWTConfigExpirationUTCField: refreshTokenExpirationTime.Unix(),
	}
	return accessTokenClaims, refreshTokenClaims
}

//GenerateFreshAccesToken generates fresh acces token
func (jwtManager JWTManager) GenerateFreshAccesToken(userIdentity string) (string, error) {
	freshAccessTokenClaims := jwtManager.generateClaimsForFreshToken(userIdentity)
	freshAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, freshAccessTokenClaims)
	freshAccessTokenString, err := freshAccessToken.SignedString(config.JWTConfigSigningSecret)
	if err != nil {
		err = errors.Wrap(err, "Error In Generating Fresh Token During Signing Fresh Token")
		return "", err
	}
	return freshAccessTokenString, err
}

//GenerateAccessAndRefreshToken will create a access and refresh token
func (jwtManager JWTManager) GenerateAccessAndRefreshToken(userIdentity string, scopesRequested []string) (string, string, error) {

	//Generating claims
	accessTokenClaims, refreshTokenClaims := jwtManager.generateClaims(userIdentity, scopesRequested)

	//Generating unsigned tokens
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	// Sign and get the complete encoded token as a string using the secret
	accessTokenString, err := accessToken.SignedString(config.JWTConfigSigningSecret)
	if err != nil {
		err = errors.Wrap(err, "Error While Signing Access Token")
		return "", "", err
	}
	refreshTokenString, err := refreshToken.SignedString(config.JWTConfigSigningSecret)
	if err != nil {
		err = errors.Wrap(err, "Error While Signing Refresh Token")
		return "", "", err
	}
	return accessTokenString, refreshTokenString, err
}

/*GenerateAccessAndRefreshTokenBasedOnRefreshToken uses and validates refresh token
to create new access and refresh token*/
func (jwtManager JWTManager) GenerateAccessAndRefreshTokenBasedOnRefreshToken(refreshTokenString string) (string, string, string, error) {
	valid, claims, msg, err := jwtManager.ValidateToken(refreshTokenString)
	if err != nil {
		err = errors.Wrap(err, "Error While Validating Refresh Token")
		return "", "", "", err
	} else if !valid {
		return "", "", config.InvalidTokenMsg, nil
	}
	userIdentity := claims[config.JWTConfigUserIdentityField].(string)
	scopes := SliceInterfaceToString(claims[config.JWTConfigScopesField].([]interface{}))
	accessToken, refreshToken, err := jwtManager.GenerateAccessAndRefreshToken(userIdentity, scopes)
	if err != nil {
		err = errors.Wrap(err, "Error While Generating Access And Refresh Token")
		return "", "", msg, err
	}
	return accessToken, refreshToken, msg, nil
}

//ValidateFreshAccessToken validates fresh acces token
func (jwtManager JWTManager) ValidateFreshAccessToken(tokenString string) (bool, string, error) {
	valid, claims, msg, err := jwtManager.ValidateToken(tokenString)
	if err != nil {
		err = errors.Wrap(err, "Error While ValidateToken Function With Fresh Token")
		return false, msg, err
	} else if !valid {
		return false, config.InvalidTokenMsg, nil
	}
	return claims[config.JWTConfigIsFreshField].(bool), "", nil
}

//ValidateToken will validate token and return its claims
func (jwtManager JWTManager) ValidateToken(tokenString string) (bool, map[string]interface{}, string, error) {
	//Parsing the access token
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return config.JWTConfigSigningSecret, nil
	})

	if err != nil {
		err = errors.Wrap(err, "Error While Parsing Token")
		return false, nil, "", err
	}

	expirationTime := timeFromFloat64(claims[config.JWTConfigExpirationUTCField].(float64))
	if !token.Valid {
		return false, nil, config.InvalidTokenMsg, nil
	} else if jwtManager.isExpired(expirationTime) {
		return false, nil, config.TokenExpiredMsg, nil
	}
	return true, claims, "", nil
}
