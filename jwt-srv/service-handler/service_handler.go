package handler

import (
	"context"

	logger "github.com/SonicRoshan/Velocity/global/logs"
	"github.com/SonicRoshan/Velocity/jwt-srv/jwt"
	proto "github.com/SonicRoshan/Velocity/jwt-srv/proto"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

//loading a logger
var log = logger.GetLogger("jwt_service.log")

//ServiceHandler is used to handle all the jwt service functions
type ServiceHandler struct {
	jwt jwt.JWT
}

//Init is used to initialize
func (serviceHandler ServiceHandler) Init() {
	serviceHandler.jwt = jwt.JWT{}
}

//FreshToken is used to generate fresh token
func (serviceHandler ServiceHandler) FreshToken(
	ctx context.Context, request *proto.JWTData, response *proto.Token) error {

	log.Debugf("Generating Fresh Token With ID %s", request.UserIdentity)
	token, err := serviceHandler.jwt.FreshToken(request.UserIdentity)
	if err != nil {
		log.Errorf("Generating Fresh Token Returned Error %+v", err)
		return errors.Wrap(err, "Error while generating fresh access token")
	}

	log.Infof("Generated Fresh Token With ID %s", request.UserIdentity)
	response.Token = token
	return nil
}

//AccessAndRefreshTokens is used to generate access and refresh token
func (serviceHandler ServiceHandler) AccessAndRefreshTokens(
	ctx context.Context, request *proto.JWTData, response *proto.AccessAndRefreshToken) error {

	log.Debugf("Generation Access And Refresh Token With ID %s", request.UserIdentity)
	accessToken, refreshToken, err := serviceHandler.jwt.AccessAndRefreshTokens(
		request.UserIdentity,
		request.Scopes,
	)
	if err != nil {
		log.Errorf("Generating Access And Refresh Token With ID %s Returned Error %+v",
			request.UserIdentity, err)
		return errors.Wrap(err, "Error while generating access and refresh token")
	}
	log.Infof("Generated Access And Refresh Token With ID %s", request.UserIdentity)
	response.AcccessToken = accessToken
	response.RefreshToken = refreshToken
	return nil
}

//RefreshTokens is used to make access and refresh token based on refresh token
func (serviceHandler ServiceHandler) RefreshTokens(
	ctx context.Context, request *proto.Token, response *proto.AccessAndRefreshToken) error {

	log.Debugf("Generating Access And Refresh Token Based On Refresh Token %s",
		request.Token)

	accessToken, refreshToken, msg, err := serviceHandler.jwt.RefreshTokens(
		request.Token)

	if err != nil {
		log.Errorf(`Generating Access And Refresh Token 
		            Based On Refresh Token %s Returned Error %+v`,
			request.Token, err)
		err = errors.Wrap(
			err, "Error while generating access and refresh token bason on refresh token")
		return err
	}

	log.Infof("Generated Access And Refresh Token Based On Refresh Token %s",
		request.Token)

	response.Message = msg
	response.AcccessToken = accessToken
	response.RefreshToken = refreshToken
	return nil
}

//ValidateToken is used to validate a token
func (serviceHandler ServiceHandler) ValidateToken(
	ctx context.Context, request *proto.Token, response *proto.Claims) error {

	log.Debugf("Validating Token %s Of Type %s", request.Token, request.TokenType)

	valid, claims, msg, err := serviceHandler.jwt.ValidateToken(request.Token, request.TokenType)
	if err != nil {
		log.Errorf("Validating Fresh Access Token %s Returned Error %+v", request.Token, err)
		return errors.Wrapf(err, "Error while validating token with type %s", request.TokenType)
	}

	response.Message = msg
	response.Valid = valid
	err = copier.Copy(&response, &claims)
	if err != nil {
		return errors.Wrap(err, "Error while copying")
	}
	log.Infof("Validated Token %s Of Type %s", request.Token, request.TokenType)
	return err
}
