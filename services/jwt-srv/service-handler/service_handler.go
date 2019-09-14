package handler

import (
	"context"

	"github.com/SonicRoshan/Velocity/global/config"
	logger "github.com/SonicRoshan/Velocity/global/logs"
	"github.com/SonicRoshan/Velocity/services/jwt-srv/jwt"
	proto "github.com/SonicRoshan/Velocity/services/jwt-srv/proto"
	"github.com/SonicRoshan/falcon"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//loading a logger
var log = logger.GetLogger("jwt_service.log")

//ServiceHandler is used to handle all the jwt service functions
type ServiceHandler struct {
	jwt        jwt.JWT
	errHandler *falcon.ErrorHandler
}

//Init is used to initialize
func (serviceHandler *ServiceHandler) Init() {
	serviceHandler.jwt = jwt.JWT{}
	serviceHandler.errHandler = falcon.NewErrorHandler()
	serviceHandler.errHandler.AddHandler(config.DefaultErrorHandler)
}

//FreshToken is used to generate fresh token
func (serviceHandler ServiceHandler) FreshToken(
	ctx context.Context, request *proto.JWTData, response *proto.Token) error {

	log.Debug("Generating Fresh Token", zap.String("ID", request.UserIdentity))

	token, err := serviceHandler.jwt.FreshToken(request.UserIdentity)
	if err != nil {
		return serviceHandler.errHandler.Check(
			err,
			"Error While Generating Fresh Access Token",
			log,
			zap.String("ID", request.UserIdentity),
		).(error)
	}

	log.Info("Generated Fresh Token", zap.String("ID", request.UserIdentity))
	response.Token = token
	return nil
}

//AccessAndRefreshTokens is used to generate access and refresh token
func (serviceHandler ServiceHandler) AccessAndRefreshTokens(
	ctx context.Context, request *proto.JWTData, response *proto.AccessAndRefreshToken) error {

	log.Debug("Generating Access And Refresh Token", zap.String("ID", request.UserIdentity))

	accessToken, refreshToken, err := serviceHandler.jwt.AccessAndRefreshTokens(
		request.UserIdentity,
		request.Scopes,
	)
	if err != nil {
		return serviceHandler.errHandler.Check(
			err,
			"Generating Access And Refresh Token Returned Error",
			log,
			zap.String("ID", request.UserIdentity),
		).(error)
	}

	log.Info("Generated Access And Refresh Token", zap.String("ID", request.UserIdentity))
	response.AcccessToken = accessToken
	response.RefreshToken = refreshToken
	return nil
}

//RefreshTokens is used to make access and refresh token based on refresh token
func (serviceHandler ServiceHandler) RefreshTokens(
	ctx context.Context, request *proto.Token, response *proto.AccessAndRefreshToken) error {

	log.Debug("Refreshing Token", zap.String("Token", request.Token))

	accessToken, refreshToken, msg, err := serviceHandler.jwt.RefreshTokens(
		request.Token)

	if err != nil {
		return serviceHandler.errHandler.Check(
			err,
			"Refreshing Token Returned Error",
			log,
			zap.String("Token", request.Token),
		).(error)
	}

	log.Info("Refreshed Token", zap.String("Token", request.Token))

	response.Message = msg
	response.AcccessToken = accessToken
	response.RefreshToken = refreshToken
	return nil
}

//ValidateToken is used to validate a token
func (serviceHandler ServiceHandler) ValidateToken(
	ctx context.Context, request *proto.Token, response *proto.Claims) error {

	log.Debug(
		"Validating Token",
		zap.String("Token", request.Token),
		zap.String("Token Type", request.TokenType),
	)

	valid, claims, msg, err := serviceHandler.jwt.ValidateToken(request.Token, request.TokenType)
	if err != nil {
		return serviceHandler.errHandler.Check(
			err,
			"Validating Token Returned Error",
			log,
			zap.String("Token", request.Token),
			zap.String("Token Type", request.TokenType),
		).(error)
	}

	response.Message = msg
	response.Valid = valid
	err = copier.Copy(&response, &claims)
	if err != nil {
		return errors.Wrap(err, "Error while copying")
	}
	log.Info(
		"Validated Token",
		zap.String("Token", request.Token),
		zap.String("Token Type", request.TokenType),
	)
	return err
}
