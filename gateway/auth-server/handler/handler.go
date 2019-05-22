package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/SonicRoshan/Velocity/global/clients"
	"github.com/SonicRoshan/Velocity/global/config"
	logger "github.com/SonicRoshan/Velocity/global/logs"
	micro "github.com/micro/go-micro"
)

var log = logger.GetLogger("auth_server.log")

//Handler will contain all the handlers
type Handler struct {
	jwtClient      clients.JWTClient
	userDataClient clients.UserDataClient
}

//getDataFromRequest gets data from url query
func (handler Handler) getDataFromRequest(w http.ResponseWriter, r *http.Request, key string) (data string) {
	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("%s Value Was Not Provided", key)
			data = ""
			handler.respond(w, nil, msg, nil)
		}
	}()
	return r.URL.Query()[key][0]
}

//Respond responds with data as json
func (handler Handler) respond(w http.ResponseWriter, data map[string]string, msg string, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Warnf("Error While Reponding %+v", err)

		w.WriteHeader(500)

		output := config.InternalServerError
		if config.AuthServerConfigShowError {
			output = err.Error()
		}

		data = map[string]string{
			config.AuthServerConfigErrField: output,
		}
	} else if msg != "" {
		data = map[string]string{
			config.AuthServerConfigErrField: msg,
		}
	}
	json.NewEncoder(w).Encode(data)
}

//Init initializes
func (handler *Handler) Init() {
	authServerService := micro.NewService(micro.Name("auth-server-srv"))
	handler.jwtClient = clients.JWTClient{}
	handler.userDataClient = clients.UserDataClient{}
	handler.jwtClient.Init(authServerService)
	handler.userDataClient.Init(authServerService)
}

//SignUpHandler signs up user and issues access and refresh token
func (handler Handler) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	username := handler.getDataFromRequest(w, r, config.DBConfigUsernameField)
	password := handler.getDataFromRequest(w, r, config.DBConfigPasswordField)
	email := handler.getDataFromRequest(w, r, config.DBConfigEmailField)
	firstName := handler.getDataFromRequest(w, r, config.DBConfigFirstNameField)
	lastName := handler.getDataFromRequest(w, r, config.DBConfigLastNameField)
	gender := handler.getDataFromRequest(w, r, config.DBConfigGenderField)
	birthdayUTC := handler.getDataFromRequest(w, r, config.DBConfigBirthdayUTCField)
	scopes := handler.getDataFromRequest(w, r, config.AuthServerConfigScopesField)

	userMainData := map[string]string{
		config.DBConfigUsernameField: username,
		config.DBConfigPasswordField: password,
		config.DBConfigEmailField:    email,
	}

	userExtraData := map[string]string{
		config.DBConfigFirstNameField:   firstName,
		config.DBConfigLastNameField:    lastName,
		config.DBConfigGenderField:      gender,
		config.DBConfigBirthdayUTCField: birthdayUTC,
	}

	userID, msg, err := handler.userDataClient.AddUser(userMainData, userExtraData)

	if msg != "" || err != nil {
		handler.respond(w, nil, msg, err)
		return
	}

	accessToken, refreshToken, msg, err := handler.jwtClient.GenerateAccessAndRefreshToken(userID, strings.Split(scopes, ","))
	output := map[string]string{
		config.AuthServerConfigAccessTokenField:  accessToken,
		config.AuthServerConfigRefreshTokenField: refreshToken,
	}
	handler.respond(w, output, msg, err)
}

//SignInHandler auths user and issues access and refresh token
func (handler Handler) SignInHandler(w http.ResponseWriter, r *http.Request) {
	username := handler.getDataFromRequest(w, r, config.DBConfigUsernameField)
	password := handler.getDataFromRequest(w, r, config.DBConfigPasswordField)
	scopes := handler.getDataFromRequest(w, r, config.AuthServerConfigScopesField)

	if username == "" || password == "" || scopes == "" {
		return
	}

	userID, msg, err := handler.userDataClient.AuthUser(username, password)
	if msg != "" || err != nil {
		handler.respond(w, nil, msg, err)
		return
	}

	accessToken, refreshToken, msg, err := handler.jwtClient.GenerateAccessAndRefreshToken(userID, strings.Split(scopes, ","))
	output := map[string]string{
		config.AuthServerConfigAccessTokenField:  accessToken,
		config.AuthServerConfigRefreshTokenField: refreshToken,
	}
	handler.respond(w, output, msg, err)
}

//SignInFreshHandler auths user and then returns a fresh access token
func (handler Handler) SignInFreshHandler(w http.ResponseWriter, r *http.Request) {
	username := handler.getDataFromRequest(w, r, config.DBConfigUsernameField)
	password := handler.getDataFromRequest(w, r, config.DBConfigPasswordField)

	if username == "" || password == "" {
		return
	}

	userID, msg, err := handler.userDataClient.AuthUser(username, password)
	if msg != "" || err != nil {
		handler.respond(w, nil, msg, err)
		return
	}

	freshAccessToken, err := handler.jwtClient.GenerateFreshAccessToken(userID)
	output := map[string]string{
		config.AuthServerConfigAccessTokenField: freshAccessToken,
	}
	handler.respond(w, output, "", err)
}

//RefreshHandler validates refresh token and issues new access and refresh token
func (handler Handler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken := handler.getDataFromRequest(w, r, config.AuthServerConfigAccessTokenField)
	if refreshToken != "" {
		return
	}

	accessToken, refreshToken, msg, err := handler.jwtClient.GenerateAccessAndRefreshTokenBasedOnRefreshToken(refreshToken)
	output := map[string]string{
		config.AuthServerConfigAccessTokenField:  accessToken,
		config.AuthServerConfigRefreshTokenField: refreshToken,
	}
	handler.respond(w, output, msg, err)
}
