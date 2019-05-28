package handler

import (
	"context"

	logger "github.com/SonicRoshan/Velocity/global/logs"
	proto "github.com/SonicRoshan/Velocity/users-srv/proto"
	"github.com/SonicRoshan/Velocity/users-srv/users"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

var log = logger.GetLogger("users_service.log")

//UsersService is to used to handle user service
type UsersService struct {
	users users.Users
}

func (usersService *UsersService) copy(toValue interface{}, fromValue interface{}) error {
	err := copier.Copy(toValue, fromValue)
	if err != nil {
		log.Errorf("Error While Copying %+v To %+v", toValue, fromValue)
		return errors.Wrap(err, "Error While Copying")
	}
	return nil
}

//Init is used to initialize
func (usersService *UsersService) Init() error {
	log.Debug("Service Initializing")
	err := usersService.users.Init()
	if err != nil {
		log.Errorf("Initializing Returned Error %+v", err)
		return errors.Wrap(err, "Error While Initializing Service")
	}
	return nil
}

//GetByUsernameOrEmail is used to handle GetByUsernameOrEmail function
func (usersService UsersService) GetByUsernameOrEmail(
	ctx context.Context,
	request *proto.GetByUsernameOrEmailRequest,
	response *proto.User) error {

	log.Debugf("Getting User By Username %s Email %s", request.Username, request.Email)
	userData, msg := usersService.users.GetByUsernameOrEmail(
		request.Username, request.Email)
	response.Message = msg
	err := usersService.copy(&userData, &response)
	if err != nil {
		return err
	}
	log.Infof("Got User By Username %s Email %s", request.Username, request.Email)
	return nil
}

//Auth is used to handle Auth function
func (usersService UsersService) Auth(
	ctx context.Context, request *proto.AuthRequest, response *proto.AuthResponse) error {

	log.Debugf("Authenticating User With Username %s Email %s", request.Username, request.Email)
	valid, userID, msg, err := usersService.users.Auth(
		request.Username, request.Email, request.Password)

	if err != nil {
		log.Errorf("Authenticating User With Username %s Email %s Returned Error %+v",
			request.Username, request.Email, err)
		return errors.Wrap(err, "Error While Authenticating User")
	}

	log.Infof("Authenticated User With Username %s Email %s", request.Username, request.Email)
	response.Message = msg
	response.Valid = valid
	response.UserID = userID
	return err
}

//Disconnect is used to disconnect
func (usersService UsersService) Disconnect() {
	usersService.Disconnect()
}
