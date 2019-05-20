package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	scopesManager "github.com/SonicRoshan/Velocity/gateway/auth-server/scopes"
	"github.com/SonicRoshan/Velocity/global/config"
	jwtSrvClient "github.com/SonicRoshan/Velocity/jwt-srv/client"
	jwtSrvProto "github.com/SonicRoshan/Velocity/jwt-srv/proto"
	userDataSrvClient "github.com/SonicRoshan/Velocity/user-data-srv/client"
	userDataSrvProto "github.com/SonicRoshan/Velocity/user-data-srv/proto"

	logger "github.com/SonicRoshan/Velocity/global/logs"
	micro "github.com/micro/go-micro"
)

//Loding Logger
var log = logger.GetLogger("auth_server.log")

//Handlers contains all http handlers for auth server
type Handlers struct {
	userDataClient userDataSrvProto.UserDataManagerService
	jwtClient      jwtSrvProto.JWTManagerService
}

//Init initializes
func (handlers *Handlers) Init() {
	authServerSrv := micro.NewService(
		micro.Name("auth-server"),
	)

	//Init of services needed
	handlers.userDataClient = userDataSrvClient.GetUserDataSrvClient(authServerSrv)

	handlers.jwtClient = jwtSrvClient.GetJWTClient(authServerSrv)
}

//SignInHandler handles sign in
func (handlers Handlers) SignInHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query()[config.AuthServerConfigUsernameField][0]
	password := r.URL.Query()[config.AuthServerConfigPasswordField][0]
	scopes := strings.Split(r.URL.Query()[config.AuthServerConfigScopesField][0], ",")

	//Validating scopes
	if !scopesManager.MatchScopesRequestedToScopesAllowed(scopes, allowedScopes) {
		json.NewEncoder(w).Encode(getError(config.InvalidScopesError.Error()))
		return
	}

	userID, response := authUser(username, password, handlers.userDataClient)
	if response != nil {
		json.NewEncoder(w).Encode(response)
		return
	}
	jwtResponse := createAccessAndRefreshToken(userID, scopes, handlers.jwtClient)
	json.NewEncoder(w).Encode(jwtResponse)
}

//SignInFreshHandler auths user and returns fresh access token
func (handlers Handlers) SignInFreshHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query()[config.AuthServerConfigUsernameField][0]
	password := r.URL.Query()[config.AuthServerConfigPasswordField][0]

	userID, response := authUser(username, password, handlers.userDataClient)
	if response != nil {
		json.NewEncoder(w).Encode(response)
		return
	}
	jwtResponse := createFreshToken(userID, handlers.jwtClient)
	json.NewEncoder(w).Encode(jwtResponse)
}
