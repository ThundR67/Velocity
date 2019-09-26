package handler

import (
	"net/http"
	"strings"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/SonicRoshan/Velocity/global/utils"
)

//SignUpHandler is used to handle sign up
func (handler Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	main, extra := handler.getUserFromURL(w, r)
	scopes := handler.getFromURL(w, r, config.AuthServerConfigScopesField)
	if scopes == "" {
		return
	}

	userID, msg := handler.users.Add(main, extra)
	if msg != "" {
		utils.GatewayRespond(w, nil, msg, nil, log)
		return
	}

	go func() {
		handler.emailVerification.SendVerification(main.Email)
	}()

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
