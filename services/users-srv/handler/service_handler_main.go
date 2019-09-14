package handler

import (
	"context"

	"github.com/SonicRoshan/Velocity/global/config"
	logger "github.com/SonicRoshan/Velocity/global/logs"
	proto "github.com/SonicRoshan/Velocity/services/users-srv/proto"
	"github.com/SonicRoshan/Velocity/services/users-srv/users"
	"github.com/SonicRoshan/falcon"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

var log = logger.GetLogger("users_service.log")

//UsersService is to used to handle user service
type UsersService struct {
	users      users.Users
	errHandler *falcon.ErrorHandler
}

func (usersService *UsersService) copy(toValue interface{}, fromValue interface{}) error {
	err := copier.Copy(toValue, fromValue)
	if err != nil {
		return usersService.errHandler.Check(
			err,
			"Error While Copying",
			zap.Any("From", fromValue),
			zap.Any("To", toValue),
		).(error)
	}
	return nil
}

//Init is used to initialize
func (usersService *UsersService) Init() error {
	log.Debug("Service Initializing")
	err := usersService.users.Init()
	if err != nil {
		return usersService.errHandler.Check(err, "Initializing Returned Error").(error)
	}
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

	err := usersService.copy(&userData, &response)
	if err != nil {
		return err
	}

	log.Info(
		"Got User By Username Or Email",
		zap.String("Username", request.Username),
		zap.String("Email", request.Email),
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

	valid, userID, msg, err := usersService.users.Auth(
		request.Username, request.Email, request.Password)

	if err != nil {
		return usersService.errHandler.Check(
			err,
			"Authenticating User Returned Error",
			zap.String("Username", request.Username),
			zap.String("Email", request.Email),
		).(error)
	}

	log.Info(
		"Authenticated User",
		zap.String("Username", request.Username),
		zap.String("Email", request.Email),
	)

	response.Message = msg
	response.Valid = valid
	response.UserID = userID
	return err
}

//Disconnect is used to disconnect
func (usersService UsersService) Disconnect() {
	usersService.Disconnect()
}
