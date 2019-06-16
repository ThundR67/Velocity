package handler

import (
	"net/http"
	"strings"

	"github.com/SonicRoshan/Velocity/global/config"
)

//SignInHandler is used to handle sign in route
func (handler Handler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := handler.getUserFromURL(w, r)
	if user.Username == "" || user.Password == "" {
		handler.respond(w, nil, "Username Or Password Not Provided", nil)
		return
	}
	scopes := handler.getFromURL(w, r, config.AuthServerConfigScopesField)
	if scopes == "" {
		return
	}

	userID, msg, err := handler.users.Auth(user.Username, user.Password)
	if err != nil || msg != "" {
		handler.respond(w, nil, msg, err)
		return
	}

	accessToken, refreshToken, msg, err := handler.jwt.AccessAndRefreshTokens(
		userID, strings.Split(scopes, ","))
	if err != nil || msg != "" {
		handler.respond(w, nil, msg, err)
		return
	}

	output := map[string]string{
		config.AuthServerConfigAccessTokenField:  accessToken,
		config.AuthServerConfigRefreshTokenField: refreshToken,
	}
	handler.respond(w, output, "", nil)
}

//SignInFreshHandler is used handle sign in fresh route where fresh access token is returned
func (handler Handler) SignInFreshHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := handler.getUserFromURL(w, r)
	if user.Username == "" || user.Password == "" {
		handler.respond(w, nil, "Username Or Password Not Provided", nil)
		return
	}

	userID, msg, err := handler.users.Auth(user.Username, user.Password)
	if err != nil || msg != "" {
		handler.respond(w, nil, msg, err)
		return
	}

	freshToken, err := handler.jwt.FreshToken(userID)
	if err != nil || msg != "" {
		handler.respond(w, nil, msg, err)
		return
	}

	output := map[string]string{
		config.AuthServerConfigAccessTokenField: freshToken,
	}
	handler.respond(w, output, "", nil)
}

//RefreshHandler is used to handle refresh route
func (handler Handler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken := handler.getFromURL(w, r, config.AuthServerConfigRefreshTokenField)
	if refreshToken == "" {
		return
	}

	accessToken, refreshToken, msg, err := handler.jwt.RefreshTokens(refreshToken)
	if err != nil || msg != "" {
		handler.respond(w, nil, msg, err)
		return
	}

	output := map[string]string{
		config.AuthServerConfigAccessTokenField:  accessToken,
		config.AuthServerConfigRefreshTokenField: refreshToken,
	}
	handler.respond(w, output, "", nil)
}
