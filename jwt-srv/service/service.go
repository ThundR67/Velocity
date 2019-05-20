package jwtservice

import (
	"context"

	"github.com/SonicRoshan/Velocity/global/config"
	jwtmanager "github.com/SonicRoshan/Velocity/jwt-srv/jwt-manager"
	proto "github.com/SonicRoshan/Velocity/jwt-srv/proto"
	logger "github.com/jex-lin/golang-logger"
)

//loading a logger
var log = logger.NewLogFile("logs/service.log")

//Service is the main jwt service struct
type Service struct {
	manager jwtmanager.JWTManager
}

//Init intitializes
func (service Service) Init() { service.manager = jwtmanager.JWTManager{} }

//GenerateFreshAccessToken generates fresh access token based on user identity
func (service Service) GenerateFreshAccessToken(ctx context.Context, request *proto.JWTData, response *proto.Token) error {
	log.Debugf("Generating Fresh Access Token With User Identity %s", request.UserIdentity)
	token, err := service.manager.GenerateFreshAccesToken(request.UserIdentity)
	if err != nil {
		log.Warnf("Generating Fresh Token With User Identity %s Returned Error %s", request.UserIdentity, err.Error())
		response.Error = err.Error()
		return err
	}
	response.Token = token
	return nil
}

//GenerateAccessAndRefreshToken generates access and refresh token based on userIdentity
func (service Service) GenerateAccessAndRefreshToken(ctx context.Context, request *proto.JWTData, response *proto.AccessAndRefreshToken) error {
	log.Debugf("Generation Access And Refresh Token For User Identity %s", request.UserIdentity)
	accessToken, refreshToken, err := service.manager.GenerateAccessAndRefreshToken(
		request.UserIdentity,
		request.Scopes,
	)
	if err != nil {
		log.Warnf("Generating Access And Refresh Token For User Identity %s Returned Error %s", request.UserIdentity, err.Error())
		response.Error = err.Error()
		return err
	}
	response.AcccessToken = accessToken
	response.RefreshToken = refreshToken
	return nil
}

//GenerateAccessAndRefreshTokenBasedOnRefreshToken generates access and refresh token based on previous refresh token
func (service Service) GenerateAccessAndRefreshTokenBasedOnRefreshToken(ctx context.Context, request *proto.Token, response *proto.AccessAndRefreshToken) error {
	log.Debugf("Generating Access And Refresh Token Based On Refresh Token %s", request.Token)
	accessToken, refreshToken, err := service.manager.GenerateAccessAndRefreshTokenBasedOnRefreshToken(request.Token)
	if err != nil {
		log.Warnf("Generating Access And Refresh Token Based On Refresh Token %s Returned Error %s", request.Token, err.Error())
		response.Error = err.Error()
		return err
	}
	response.AcccessToken = accessToken
	response.RefreshToken = refreshToken
	return nil
}

//ValidateFreshAccessToken validates fresh access token
func (service Service) ValidateFreshAccessToken(ctx context.Context, request *proto.Token, response *proto.Claims) error {
	log.Debugf("Validating Fresh Access Token %s", request.Token)
	valid, err := service.manager.ValidateFreshAccessToken(request.Token)
	if err != nil {
		log.Warnf("Validating Fresh Access Token %s Returned Error %s", request.Token, err.Error())
		response.Error = err.Error()
		return err
	}
	response.Valid = valid
	return err
}

//ValidateToken validates access token and returnes its claims
func (service Service) ValidateToken(ctx context.Context, request *proto.Token, response *proto.Claims) error {
	log.Debugf("Validating Token %s", request.Token)
	valid, claims, err := service.manager.ValidateToken(request.Token)
	if err != nil {
		log.Warnf("Validating Token %s Returned Error %s", request.Token, err.Error())
		response.Error = err.Error()
		return err
	}
	response.Valid = valid
	if !valid {
		return nil
	}

	response.JwtData.UserIdentity = claims[config.JWTConfigUserIdentityField].(string)
	response.JwtData.Scopes = jwtmanager.SliceInterfaceToString(claims[config.JWTConfigScopesField].([]interface{}))
	response.CreationUTC = claims[config.JWTConfigCreationUTCField].(float64)
	response.ExpirationUTC = claims[config.JWTConfigExpirationUTCField].(float64)
	return err
}
