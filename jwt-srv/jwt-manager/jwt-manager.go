package jwtmanager

import (
	"fmt"
	"time"

	"github.com/SonicRoshan/Velocity/global/config"
	jwt "github.com/dgrijalva/jwt-go"
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

type customClaim struct {
	UserIdentity  string   `json:"userIdentity"`
	Fresh         bool     `json:"fresh"`
	Scopes        []string `json:"scopes"`
	CreationUTC   float64  `json:"creationUTC"`
	ExpirationUTC float64  `json:"expirationUTC"`
	jwt.Claims
}

func (cc customClaim) toMap() map[string]interface{} {
	return map[string]interface{}{
		config.JWTConfigUserIdentityField:  cc.UserIdentity,
		config.JWTConfigIsFreshField:       cc.Fresh,
		config.JWTConfigScopesField:        cc.Scopes,
		config.JWTConfigCreationUTCField:   cc.CreationUTC,
		config.JWTConfigExpirationUTCField: cc.ExpirationUTC,
	}
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
		return "", "", err
	}
	refreshTokenString, err := refreshToken.SignedString(config.JWTConfigSigningSecret)
	return accessTokenString, refreshTokenString, err
}

/*GenerateAccessAndRefreshTokenBasedOnRefreshToken uses and validates refresh token
to create new access and refresh token*/
func (jwtManager JWTManager) GenerateAccessAndRefreshTokenBasedOnRefreshToken(refreshTokenString string) (string, string, error) {
	valid, claims, err := jwtManager.ValidateToken(refreshTokenString)
	if !valid || err != nil {
		return "", "", config.InvalidTokenError
	}
	userIdentity := claims[config.JWTConfigUserIdentityField].(string)
	scopes := SliceInterfaceToString(claims[config.JWTConfigScopesField].([]interface{}))
	return jwtManager.GenerateAccessAndRefreshToken(userIdentity, scopes)
}

//ValidateFreshAccessToken validates fresh acces token
func (jwtManager JWTManager) ValidateFreshAccessToken(tokenString string) (bool, error) {
	valid, claims, err := jwtManager.ValidateToken(tokenString)
	if err != nil {
		return false, err
	} else if !valid {
		return false, config.InvalidTokenError
	}
	return claims[config.JWTConfigIsFreshField].(bool), nil
}

//ValidateToken will validate token and return its claims
func (jwtManager JWTManager) ValidateToken(tokenString string) (bool, map[string]interface{}, error) {
	//Parsing the access token
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return config.JWTConfigSigningSecret, nil
	})

	if err != nil {
		return false, nil, err
	}

	expirationTime := timeFromFloat64(claims[config.JWTConfigExpirationUTCField].(float64))
	if !token.Valid {
		return false, nil, config.InvalidTokenError
	} else if jwtManager.isExpired(expirationTime) {
		return false, nil, config.TokenExpiredError
	}
	return true, claims, nil
}
