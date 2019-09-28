package handler

import (
	"context"

	"github.com/SonicRoshan/falcon"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/SonicRoshan/Velocity/global/logger"
	"github.com/SonicRoshan/Velocity/services/jwt-srv/jwt"
	proto "github.com/SonicRoshan/Velocity/services/jwt-srv/proto"
)

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

	token := serviceHandler.jwt.FreshToken(request.UserIdentity)

	log.Info("Generated Fresh Token", zap.String("ID", request.UserIdentity))
	response.Token = token

	return nil
}

//AccessAndRefreshTokens is used to generate access and refresh tokens
func (serviceHandler ServiceHandler) AccessAndRefreshTokens(
	ctx context.Context,
	request *proto.JWTData,
	response *proto.AccessAndRefreshToken) error {

	log.Debug("Generating Access And Refresh Token",
		zap.String("ID", request.UserIdentity))

	accessToken, refreshToken := serviceHandler.jwt.AccessAndRefreshTokens(
		request.UserIdentity,
		request.Scopes,
	)

	log.Info("Generated Access And Refresh Token",
		zap.String("ID", request.UserIdentity))

	response.AcccessToken = accessToken
	response.RefreshToken = refreshToken

	return nil
}

//RefreshTokens is used to make access and refresh token based on refresh token
func (serviceHandler ServiceHandler) RefreshTokens(
	ctx context.Context,
	request *proto.Token,
	response *proto.AccessAndRefreshToken) error {

	log.Debug("Refreshing Token", zap.String("Token", request.Token))

	accessToken, refreshToken, msg, err := serviceHandler.jwt.RefreshTokens(
		request.Token,
	)

	if err != nil {
		return serviceHandler.errHandler.Check(
			err,
			"Refreshing Token Returned Error",
			log,
			zap.String("Token", request.Token),
		).(error)
	}

	log.Info(
		"Refreshed Tokens",
		zap.String("Token", request.Token),
		zap.String("AccesToken", accessToken),
		zap.String("RefreshToken", refreshToken),
	)

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

	valid, claims, msg, err := serviceHandler.jwt.ValidateToken(
		request.Token,
		request.TokenType,
	)

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
	copier.Copy(&response, &claims)

	log.Info(
		"Validated Token",
		zap.String("Token", request.Token),
		zap.String("Token Type", request.TokenType),
	)

	return err
}
