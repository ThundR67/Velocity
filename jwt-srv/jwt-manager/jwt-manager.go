package jwtmanager

import (
	"fmt"
	"time"

	config "github.com/SonicRoshan/Velocity/jwt-srv/config"
	"github.com/SonicRoshan/Velocity/jwt-srv/jwt-manager/scopes"
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
		config.ConfigUserIdentityField:  cc.UserIdentity,
		config.ConfigIsFreshField:       cc.Fresh,
		config.ConfigScopesField:        cc.Scopes,
		config.ConfigCreationUTCField:   cc.CreationUTC,
		config.ConfigExpirationUTCField: cc.ExpirationUTC,
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
	freshAccessTokenExpirationTime := currentTime.Add(time.Minute * config.ConfigFreshAccessTokenExpirationTimeMinutes)
	return jwt.MapClaims{
		config.ConfigUserIdentityField:  userIdentity,
		config.ConfigIsFreshField:       true,
		config.ConfigCreationUTCField:   currentTime.Unix(),
		config.ConfigExpirationUTCField: freshAccessTokenExpirationTime.Unix(),
	}
}

//generateClaims for acces and refresh token
func (jwtManager JWTManager) generateClaims(userIdentity string, scopes []string) (jwt.MapClaims, jwt.MapClaims) {
	currentTime := time.Now()
	accessTokenExpirationTime := currentTime.Add(time.Minute * config.ConfigAccessTokenExpirationTimeMinutes)
	refreshTokenExpirationTime := currentTime.Add(time.Hour * 24 * config.ConfigRefreshTokenExpirationTimeDays)
	accessTokenClaims := jwt.MapClaims{
		config.ConfigUserIdentityField:  userIdentity,
		config.ConfigIsFreshField:       false,
		config.ConfigScopesField:        scopes,
		config.ConfigCreationUTCField:   currentTime.Unix(),
		config.ConfigExpirationUTCField: accessTokenExpirationTime.Unix(),
	}

	refreshTokenClaims := jwt.MapClaims{
		config.ConfigUserIdentityField:  userIdentity,
		config.ConfigScopesField:        scopes,
		config.ConfigCreationUTCField:   currentTime,
		config.ConfigExpirationUTCField: refreshTokenExpirationTime.Unix(),
	}
	return accessTokenClaims, refreshTokenClaims
}

//GenerateFreshAccesToken generates fresh acces token
func (jwtManager JWTManager) GenerateFreshAccesToken(userIdentity string) (string, error) {
	freshAccessTokenClaims := jwtManager.generateClaimsForFreshToken(userIdentity)
	freshAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, freshAccessTokenClaims)
	freshAccessTokenString, err := freshAccessToken.SignedString(config.ConfigSigningSecret)
	return freshAccessTokenString, err
}

//GenerateAccessAndRefreshToken will create a access and refresh token
func (jwtManager JWTManager) GenerateAccessAndRefreshToken(userIdentity string, scopesRequested []string) (string, string, error) {
	//Check if scopes are allowed
	if !scopes.MatchScopesRequestedToScopesAllowed(scopesRequested, allowedScopes) {
		return "", "", InvalidScopesError{}
	}

	//Generating claims
	accessTokenClaims, refreshTokenClaims := jwtManager.generateClaims(userIdentity, scopesRequested)

	//Generating unsigned tokens
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	// Sign and get the complete encoded token as a string using the secret
	accessTokenString, err := accessToken.SignedString(config.ConfigSigningSecret)
	if err != nil {
		return "", "", err
	}
	refreshTokenString, err := refreshToken.SignedString(config.ConfigSigningSecret)
	return accessTokenString, refreshTokenString, err
}

/*GenerateAccessAndRefreshTokenBasedOnRefreshToken uses and validates refresh token
to create new access and refresh token*/
func (jwtManager JWTManager) GenerateAccessAndRefreshTokenBasedOnRefreshToken(refreshTokenString string) (string, string, error) {
	valid, claims, err := jwtManager.ValidateToken(refreshTokenString)
	if !valid || err != nil {
		return "", "", InvalidTokenError{}
	}
	userIdentity := claims[config.ConfigUserIdentityField].(string)
	scopes := SliceInterfaceToString(claims[config.ConfigScopesField].([]interface{}))
	return jwtManager.GenerateAccessAndRefreshToken(userIdentity, scopes)
}

//ValidateFreshAccessToken validates fresh acces token
func (jwtManager JWTManager) ValidateFreshAccessToken(tokenString string) (bool, error) {
	valid, claims, err := jwtManager.ValidateToken(tokenString)
	if err != nil {
		return false, err
	} else if !valid {
		return false, InvalidTokenError{}
	}
	return claims[config.ConfigIsFreshField].(bool), nil
}

//ValidateToken will validate token and return its claims
func (jwtManager JWTManager) ValidateToken(tokenString string) (bool, map[string]interface{}, error) {
	//Parsing the access token
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return config.ConfigSigningSecret, nil
	})

	if err != nil {
		return false, nil, err
	}

	expirationTime := timeFromFloat64(claims[config.ConfigExpirationUTCField].(float64))
	if !token.Valid {
		return false, nil, InvalidTokenError{}
	} else if jwtManager.isExpired(expirationTime) {
		return false, nil, TokenExpiredError{}
	}
	return true, claims, nil
}
