package handler

import (
	"net/http"
	"strings"

	"github.com/SonicRoshan/Velocity/global/config"
)

//splitUser is used to split user data into main data and extra data
func splitUser(user config.UserType) (config.UserType, config.UserType) {
	mainData := config.UserType{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
	extraData := config.UserType{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		BirthdayUTC: user.BirthdayUTC,
		Gender:      user.Gender,
	}
	return mainData, extraData
}

//SignUpHandler is used to handle sign up
func (handler Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	user := handler.getUserFromURL(w, r)
	scopes := handler.getFromURL(w, r, config.AuthServerConfigScopesField)
	if scopes == "" {
		return
	}
	main, extra := splitUser(user)

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
