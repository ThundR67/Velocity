package handler

import (
	"net/http"
	"strings"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/SonicRoshan/Velocity/global/utils"
)

//SignInHandler is used to handle sign in route
func (handler Handler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := handler.getUserFromURL(w, r)
	if user.Username == "" || user.Password == "" {
		utils.GatewayRespond(w, nil, "Username Or Password Not Provided", nil, log)
		return
	}
	scopes := handler.getFromURL(w, r, config.AuthServerConfigScopesField)
	if scopes == "" {
		return
	}

	userID, msg, err := handler.users.Auth(user.Username, user.Password)
	if err != nil || msg != "" {
		utils.GatewayRespond(w, nil, msg, err, log)
		return
	}

	accessToken, refreshToken, msg, err := handler.jwt.AccessAndRefreshTokens(
		userID, strings.Split(scopes, ","))
	if err != nil || msg != "" {
		utils.GatewayRespond(w, nil, msg, err, log)
		return
	}

	output := map[string]string{
		config.AuthServerConfigAccessTokenField:  accessToken,
		config.AuthServerConfigRefreshTokenField: refreshToken,
	}
	utils.GatewayRespond(w, output, "", nil, log)
}

//SignInFreshHandler is used handle sign in fresh route where fresh access token is returned
func (handler Handler) SignInFreshHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := handler.getUserFromURL(w, r)
	if user.Username == "" || user.Password == "" {
		utils.GatewayRespond(w, nil, "Username Or Password Not Provided", nil, log)
		return
	}

	userID, msg, err := handler.users.Auth(user.Username, user.Password)
	if err != nil || msg != "" {
		utils.GatewayRespond(w, nil, msg, err, log)
		return
	}

	freshToken := handler.jwt.FreshToken(userID)
	if msg != "" {
		utils.GatewayRespond(w, nil, msg, nil, log)
		return
	}

	output := map[string]string{
		config.AuthServerConfigAccessTokenField: freshToken,
	}
	utils.GatewayRespond(w, output, "", nil, log)
}

//RefreshHandler is used to handle refresh route
func (handler Handler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken := handler.getFromURL(w, r, config.AuthServerConfigRefreshTokenField)
	if refreshToken == "" {
		return
	}

	accessToken, refreshToken, msg, err := handler.jwt.RefreshTokens(refreshToken)
	if err != nil || msg != "" {
		utils.GatewayRespond(w, nil, msg, err, log)
		return
	}

	output := map[string]string{
		config.AuthServerConfigAccessTokenField:  accessToken,
		config.AuthServerConfigRefreshTokenField: refreshToken,
	}
	utils.GatewayRespond(w, output, "", nil, log)
}
