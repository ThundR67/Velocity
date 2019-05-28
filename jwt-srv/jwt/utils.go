package jwt

import (
	"fmt"
	"time"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/dgrijalva/jwt-go"
)

//makeClaims is used to generate claims for a token
func makeClaims(
	userIdentity string,
	scopes []string,
	isFresh bool,
	isRefresh bool,
	expiration time.Duration) config.JWTClaims {

	currentTime := time.Now()
	expirationTime := currentTime.Add(expiration)

	return config.JWTClaims{
		UserIdentity:  userIdentity,
		IsFresh:       isFresh,
		IsRefresh:     isRefresh,
		Scopes:        scopes,
		CreationUTC:   currentTime.Unix(),
		ExpirationUTC: expirationTime.Unix(),
	}
}

//freshTokenClaims is used to make claims for a fresh access token
func freshTokenClaims(userIdentity string) config.JWTClaims {
	return makeClaims(userIdentity, nil, true, false, config.JWTFreshAccessExpirationMinutes)
}

//accessTokenClaims is used to make claims for access token
func accessTokenClaims(userIdentity string, scopes []string) config.JWTClaims {
	return makeClaims(userIdentity, scopes, false, false, config.JWTAccessExpirationMinutes)
}

//refreshTokenClaims is used to make claims for refresh token
func refreshTokenClaims(userIdentity string, scopes []string) config.JWTClaims {
	return makeClaims(userIdentity, scopes, false, true, config.JWTRefreshExpirationDays)
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return config.JWTSecret, nil
}
