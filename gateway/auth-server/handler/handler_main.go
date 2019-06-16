package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SonicRoshan/Velocity/global/clients"
	"github.com/SonicRoshan/Velocity/global/config"
	logger "github.com/SonicRoshan/Velocity/global/logs"
	"github.com/gorilla/schema"
	micro "github.com/micro/go-micro"
)

var log = logger.GetLogger("auth_server.log")
var decoder = schema.NewDecoder()

//Handler is used to handle all routes of auth server
type Handler struct {
	jwt   clients.JWTClient
	users clients.UsersClient
}

//getFromURL is used to get data from url query
func (handler Handler) getFromURL(
	w http.ResponseWriter, r *http.Request, key string) (data string) {

	defer func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("%s Value Was Not Provided", key)
			data = ""
			handler.respond(w, nil, msg, nil)
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
		handler.respond(w, nil, "", err)
		return
	}
	err = decoder.Decode(&extraData, r.URL.Query())
	if err != nil {
		handler.respond(w, nil, "", err)
		return
	}
	return
}

//Respond is used to respond to client with json
func (handler Handler) respond(
	w http.ResponseWriter, data map[string]string, msg string, err error) {

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		log.Errorf("Error While Reponding %+v", err)

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

//Init is used to initialise
func (handler *Handler) Init() {
	authSrv := micro.NewService(micro.Name(config.AuthServerService))
	handler.jwt = clients.NewJWTClient(authSrv)
	handler.users = clients.NewUsersClient(authSrv)
	decoder.IgnoreUnknownKeys(true)
}
