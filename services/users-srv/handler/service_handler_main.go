package handler

import (
	"context"

	"github.com/SonicRoshan/falcon"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"

	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/SonicRoshan/Velocity/global/logger"
	proto "github.com/SonicRoshan/Velocity/services/users-srv/proto"
	"github.com/SonicRoshan/Velocity/services/users-srv/users"
)

var log = logger.GetLogger("users_service.log")

//UsersService is to used to handle user service
type UsersService struct {
	users      users.Users
	errHandler *falcon.ErrorHandler
}

//Init is used to initialize
func (usersService *UsersService) Init() error {

	log.Debug("Service Initializing")
	usersService.users.Init()

	usersService.errHandler = falcon.NewErrorHandler()
	usersService.errHandler.AddHandler(config.DefaultErrorHandler)

	return nil
}

//GetByUsernameOrEmail is used to handle GetByUsernameOrEmail function
func (usersService UsersService) GetByUsernameOrEmail(
	ctx context.Context,
	request *proto.GetByUsernameOrEmailRequest,
	response *proto.UserMain) error {

	log.Debug(
		"Getting User By Username Or Email",
		zap.String("Username", request.Username),
		zap.String("Email", request.Email),
	)

	userData, msg := usersService.users.GetByUsernameOrEmail(
		request.Username, request.Email)
	response.Message = msg

	log.Info("Got Data From Low Level",
		zap.Any("UserData", userData),
		zap.String("message", msg),
	)

	copier.Copy(&response, &userData)

	log.Info(
		"Got User By Username Or Email",
		zap.String("Username", request.Username),
		zap.String("Email", request.Email),
		zap.Any("UserData", userData),
		zap.Any("Response", response),
	)

	return nil
}

//Auth is used to handle Auth function
func (usersService UsersService) Auth(
	ctx context.Context, request *proto.AuthRequest, response *proto.AuthResponse) error {

	log.Debug(
		"Authenticating User",
		zap.String("Username", request.Username),
		zap.String("Email", request.Email),
	)

	valid, userID, msg := usersService.users.Auth(
		request.Username, request.Email, request.Password)

	log.Info(
		"Authenticated User",
		zap.String("Username", request.Username),
		zap.String("Email", request.Email),
	)

	response.Message = msg
	response.Valid = valid
	response.UserID = userID
	return nil
}

//Disconnect is used to disconnect
func (usersService UsersService) Disconnect() {
	usersService.Disconnect()
}
