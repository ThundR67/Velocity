package handler

import (
	"net/http"
	"strings"

	"github.com/SonicRoshan/Velocity/global/config"
)

//SignUpHandler is used to handle sign up
func (handler Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	main, extra := handler.getUserFromURL(w, r)
	scopes := handler.getFromURL(w, r, config.AuthServerConfigScopesField)
	if scopes == "" {
		return
	}

	userID, msg, err := handler.users.Add(main, extra)
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
