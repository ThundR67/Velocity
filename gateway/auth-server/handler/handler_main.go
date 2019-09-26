package handler

import (
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"github.com/SonicRoshan/Velocity/global/clients"
	"github.com/SonicRoshan/Velocity/global/config"
	logger "github.com/SonicRoshan/Velocity/global/logs"
	"github.com/SonicRoshan/Velocity/global/utils"
	"github.com/gorilla/schema"
	micro "github.com/micro/go-micro"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var log = logger.GetLogger("auth_server.log")
var decoder = schema.NewDecoder()

//Handler is used to handle all routes of auth server
type Handler struct {
	jwt               clients.JWTClient
	users             clients.UsersClient
	emailVerification clients.EmailVerificationClient
}

//getFromURL is used to get data from url query
func (handler Handler) getFromURL(
	w http.ResponseWriter, r *http.Request, key string) (data string) {

	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("%s Value Was Not Provided", key)
			data = ""
			utils.GatewayRespond(w, nil, msg, nil, log)
		}
	}()

	return r.URL.Query()[key][0]
}

//getUserFromUrl is used to get user main and extra data from url
func (handler Handler) getUserFromURL(
	w http.ResponseWriter,
	r *http.Request) (mainData config.UserMain, extraData config.UserExtra) {

	err := decoder.Decode(&mainData, r.URL.Query())
	if err != nil {
		utils.GatewayRespond(w, nil, "", err, log)
		return
	}
	err = decoder.Decode(&extraData, r.URL.Query())
	if err != nil {
		utils.GatewayRespond(w, nil, "", err, log)
		return
	}
	return
}

//Init is used to initialise
func (handler *Handler) Init() {
	authSrv := micro.NewService(micro.Name(config.AuthServerService))
	handler.jwt = clients.NewJWTClient(authSrv)
	handler.users = clients.NewUsersClient(authSrv)
	handler.emailVerification = clients.NewEmailVerificationClient(authSrv)
	decoder.IgnoreUnknownKeys(true)
}
