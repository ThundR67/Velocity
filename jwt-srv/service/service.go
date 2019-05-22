package jwtservice

import (
	"context"

	"github.com/SonicRoshan/Velocity/global/config"
	logger "github.com/SonicRoshan/Velocity/global/logs"
	jwtmanager "github.com/SonicRoshan/Velocity/jwt-srv/jwt-manager"
	proto "github.com/SonicRoshan/Velocity/jwt-srv/proto"
)

//loading a logger
var log = logger.GetLogger("jwt_service.log")

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
		log.Warnf("Generating Fresh Token With User Identity %s Returned Error %+v", request.UserIdentity, err)
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
		log.Warnf("Generating Access And Refresh Token For User Identity %s Returned Error %+v", request.UserIdentity, err)
		return err
	}
	response.AcccessToken = accessToken
	response.RefreshToken = refreshToken
	return nil
}

//GenerateAccessAndRefreshTokenBasedOnRefreshToken generates access and refresh token based on previous refresh token
func (service Service) GenerateAccessAndRefreshTokenBasedOnRefreshToken(ctx context.Context, request *proto.Token, response *proto.AccessAndRefreshToken) error {
	log.Debugf("Generating Access And Refresh Token Based On Refresh Token %s", request.Token)
	accessToken, refreshToken, msg, err := service.manager.GenerateAccessAndRefreshTokenBasedOnRefreshToken(request.Token)
	if err != nil {
		log.Warnf("Generating Access And Refresh Token Based On Refresh Token %s Returned Error %+v", request.Token, err)
		return err
	}
	response.Message = msg
	response.AcccessToken = accessToken
	response.RefreshToken = refreshToken
	return nil
}

//ValidateFreshAccessToken validates fresh access token
func (service Service) ValidateFreshAccessToken(ctx context.Context, request *proto.Token, response *proto.Claims) error {
	log.Debugf("Validating Fresh Access Token %s", request.Token)
	valid, msg, err := service.manager.ValidateFreshAccessToken(request.Token)
	if err != nil {
		log.Warnf("Validating Fresh Access Token %s Returned Error %+v", request.Token, err)
		return err
	}
	response.Message = msg
	response.Valid = valid
	return err
}

//ValidateToken validates access token and returnes its claims
func (service Service) ValidateToken(ctx context.Context, request *proto.Token, response *proto.Claims) error {
	log.Debugf("Validating Token %s", request.Token)
	valid, claims, msg, err := service.manager.ValidateToken(request.Token)
	if err != nil {
		log.Warnf("Validating Token %s Returned Error %+v", request.Token, err)
		return err
	}
	response.Valid = valid
	if !valid {
		return nil
	}

	response.Message = msg
	response.JwtData.UserIdentity = claims[config.JWTConfigUserIdentityField].(string)
	response.JwtData.Scopes = jwtmanager.SliceInterfaceToString(claims[config.JWTConfigScopesField].([]interface{}))
	response.CreationUTC = claims[config.JWTConfigCreationUTCField].(float64)
	response.ExpirationUTC = claims[config.JWTConfigExpirationUTCField].(float64)
	return err
}
