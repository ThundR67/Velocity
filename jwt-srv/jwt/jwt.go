package jwt

import (
	"time"

	"github.com/SonicRoshan/Velocity/global/config"
	goJwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

//JWT is a low level jason web token manager
type JWT struct{}

//isExpired checks if a time.Duration is expired
func (jwt JWT) isExpired(expirationTime time.Time) bool {
	return time.Now().After(expirationTime)
}

//FreshToken is used to generate a fresh access token
func (jwt JWT) FreshToken(userIdentity string) (string, error) {
	freshTokenClaims := freshTokenClaims(userIdentity)
	freshToken := goJwt.NewWithClaims(goJwt.SigningMethodHS256, freshTokenClaims)
	freshTokenString, err := freshToken.SignedString(config.JWTSecret)
	if err != nil {
		err = errors.Wrap(err, "Error In Generating Fresh Token During Signing Fresh Token")
		return "", err
	}
	return freshTokenString, nil
}

//AccessAndRefreshTokens is used to create access and refresh token
func (jwt JWT) AccessAndRefreshTokens(
	userIdentity string, scopesRequested []string) (string, string, error) {

	accessClaims := accessTokenClaims(userIdentity, scopesRequested)
	refreshClaims := refreshTokenClaims(userIdentity, scopesRequested)

	accessToken := goJwt.NewWithClaims(goJwt.SigningMethodHS256, accessClaims)
	refreshToken := goJwt.NewWithClaims(goJwt.SigningMethodHS256, refreshClaims)

	accessTokenString, err := accessToken.SignedString(config.JWTSecret)
	if err != nil {
		err = errors.Wrap(err, "Error While Signing Access Token")
		return "", "", err
	}
	refreshTokenString, err := refreshToken.SignedString(config.JWTSecret)
	if err != nil {
		err = errors.Wrap(err, "Error While Signing Refresh Token")
		return "", "", err
	}
	return accessTokenString, refreshTokenString, nil
}

//RefreshTokens is used to generate new access and refresh token based on refresh token
func (jwt JWT) RefreshTokens(refreshTokenString string) (string, string, string, error) {
	valid, claims, msg, err := jwt.ValidateToken(refreshTokenString, config.TokenTypeRefresh)

	if err != nil {
		err = errors.Wrap(err, "Error While Validating Refresh Token")
		return "", "", "", err
	} else if !valid {
		return "", "", config.InvalidTokenMsg, nil
	}

	userIdentity := claims.UserIdentity
	scopes := claims.Scopes
	accessToken, refreshToken, err := jwt.AccessAndRefreshTokens(userIdentity, scopes)
	if err != nil {
		err = errors.Wrap(err, "Error While Generating Access And Refresh Token")
		return "", "", msg, err
	}
	return accessToken, refreshToken, msg, nil
}

//ValidateToken is used to validate a token
func (jwt JWT) ValidateToken(tokenString, tokenType string) (bool, config.JWTClaims, string, error) {
	var claims config.JWTClaims
	token, err := goJwt.ParseWithClaims(tokenString, &claims, jwtKeyFunc)

	if err != nil {
		err = errors.Wrap(err, "Error While Parsing Token")
		return false, config.JWTClaims{}, "", err
	}

	expirationTime := time.Unix(claims.ExpirationUTC, 0)
	if !token.Valid {
		return false, config.JWTClaims{}, config.InvalidTokenMsg, nil
	} else if jwt.isExpired(expirationTime) {
		return false, config.JWTClaims{}, config.TokenExpiredMsg, nil
	}

	switch tokenType {
	case config.TokenTypeAccess:
		return !claims.IsFresh && !claims.IsRefresh, claims, "", nil
	case config.TokenTypeFresh:
		return claims.IsFresh, claims, "", nil
	case config.TokenTypeRefresh:
		return claims.IsRefresh, claims, "", nil
	}

	return false, config.JWTClaims{}, "", errors.New("Invalid Token Type")
}
