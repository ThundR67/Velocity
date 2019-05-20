package handlers

import (
	"context"

	"github.com/SonicRoshan/Velocity/global/config"
	jwtSrvProto "github.com/SonicRoshan/Velocity/jwt-srv/proto"
	userDataSrvProto "github.com/SonicRoshan/Velocity/user-data-srv/proto"
)

//Helper functions

func getError(msg string) map[string]string {
	return map[string]string{
		config.AuthServerConfigErrField: msg,
	}
}

//authUser talks to services and auths user
func authUser(username, password string, client userDataSrvProto.UserDataManagerService) (string, map[string]string) {
	request := userDataSrvProto.AuthUserRequest{
		Username: username,
		Password: password,
	}
	response, err := client.AuthUser(context.TODO(), &request)
	if err != nil {
		if err.Error() == config.InvalidUsernameOrEmailError.Error() || err.Error() == config.InvalidPasswordError.Error() {
			return "", getError(err.Error())
		}
		log.Criticalf("Internal Server Error While Authenticating User %s", err.Error())
		return "", getError(config.InternalServerError)
	}

	return response.UserID, nil
}

//createAccessAndRefreshToken talks to services and creates access and refresh token
func createAccessAndRefreshToken(userIdentity string, scopes []string, client jwtSrvProto.JWTManagerService) map[string]string {
	jwtRequest := jwtSrvProto.JWTData{
		UserIdentity: userIdentity,
		Scopes:       scopes,
	}
	jwtResponse, err := client.GenerateAccessAndRefreshToken(context.TODO(), &jwtRequest)
	if err != nil {
		log.Critical("Internal Server Error While Creating Fresh Token %s", err.Error())
		return getError(config.InternalServerError)
	}

	return map[string]string{
		config.AuthServerConfigAccessTokenField:  jwtResponse.AcccessToken,
		config.AuthServerConfigRefreshTokenField: jwtResponse.RefreshToken,
	}
}

//createFreshToken talks to services and creates fresh acces token
func createFreshToken(userIdentity string, client jwtSrvProto.JWTManagerService) map[string]string {
	jwtRequest := jwtSrvProto.JWTData{
		UserIdentity: userIdentity,
	}
	jwtResponse, err := client.GenerateFreshAccessToken(context.TODO(), &jwtRequest)
	if err != nil {
		log.Critical("Internal Server Error While Creating Fresh Token %s", err.Error())
		return getError(config.InternalServerError)
	}

	return map[string]string{
		config.AuthServerConfigAccessTokenField: jwtResponse.Token,
	}
}
