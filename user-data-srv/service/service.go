package userdataservice

import (
	"context"

	"github.com/SonicRoshan/Velocity/global/config"
	logger "github.com/SonicRoshan/Velocity/global/logs"
	proto "github.com/SonicRoshan/Velocity/user-data-srv/proto"
	userDataManager "github.com/SonicRoshan/Velocity/user-data-srv/user-data-manager"
)

//Loading logger
var serviceLog = logger.GetLogger("user_data_service_log.log")

//castFromStringToInterface converts map[string]string to map[string]interface{}
func castFromStringToInterface(mapInput map[string]string) map[string]interface{} {
	mapOutput := make(map[string]interface{})
	for key, value := range mapInput {
		mapOutput[key] = value
	}
	return mapOutput
}

//castFromInterfaceToString Converts map[string]interface{} to map[string]string
func castFromInterfaceToString(mapInput map[string]interface{}) map[string]string {
	mapOutput := make(map[string]string)
	for key, value := range mapInput {
		mapOutput[key] = value.(string)
	}
	return mapOutput
}

//UserDataService Is The Main Service Struct
type UserDataService struct {
	userDataClient userDataManager.UserDataManager
}

//Init Initializes
func (userDataService *UserDataService) Init() error {
	err := userDataService.userDataClient.Init()
	return err
}

//AddUser Handles AddUser Function
func (userDataService UserDataService) AddUser(ctx context.Context, request *proto.UserData, response *proto.AddUserResponse) error {
	userID, err := userDataService.userDataClient.AddUser(castFromStringToInterface(request.Data))
	if err != nil {
		if err != config.UsernameExistError && err != config.EmailExistError {
			serviceLog.Criticalf("AddUser With Data %v Returned Error %s", request.Data, err)
		}
		response.Error = err.Error()
		return err
	}
	response.UserID = userID
	return err
}

//GetUser Handles GetUser Function
func (userDataService UserDataService) GetUser(ctx context.Context, request *proto.GetUserRequest, response *proto.UserData) error {
	userData, err := userDataService.userDataClient.GetUser(request.UserID)
	if err != nil {
		serviceLog.Criticalf("GetUser With UserID %s Returned Error %s", request.UserID, err)
		response.Error = err.Error()
	}
	response.Data = castFromInterfaceToString(userData)
	return err
}

//GetUserByUsernameOrEmail Handles GetUserByUsernameOrEmail Function
func (userDataService UserDataService) GetUserByUsernameOrEmail(ctx context.Context, request *proto.GetUserByUsernameOrEmailRequest, response *proto.UserData) error {
	userData, err := userDataService.userDataClient.GetUserByUsernameOrEmail(request.Username, request.Email)
	if err != nil {
		serviceLog.Criticalf("GetUserByUsernameOrEmail With Username %s Email %s Returned Error %s", request.Username, request.Email, err)
		response.Error = err.Error()
	}
	response.Data = castFromInterfaceToString(userData)
	return err
}

//AuthUser Handles AuthUser Function
func (userDataService UserDataService) AuthUser(ctx context.Context, request *proto.AuthUserRequest, response *proto.AuthUserResponse) error {
	valid, userID, err := userDataService.userDataClient.AuthUser(request.Username, request.Email, request.Password)
	if err != nil {
		if err != config.InvalidUsernameOrEmailError && err != config.InvalidPasswordError {
			serviceLog.Criticalf("AuthUser With Username %s Email %s Returned Error %s", request.Username, request.Email, err)
		}
		response.Error = err.Error()
		return err
	}
	response.Valid = valid
	response.UserID = userID
	return err
}

//UpdateUser Handles UpdateUser Function
func (userDataService UserDataService) UpdateUser(ctx context.Context, request *proto.UpdateUserRequest, response *proto.UpdateUserResponse) error {
	err := userDataService.userDataClient.UpdateUser(request.UserID, request.Field, request.UpdatedValue)
	if err != nil {
		serviceLog.Criticalf("UpdateUser With Id %s Returned Error %s", request.UserID, err)
		response.Error = err.Error()
	}
	return err

}

//DeleteUser Handles DeleteUser Function
func (userDataService UserDataService) DeleteUser(ctx context.Context, request *proto.DeleteUserRequest, response *proto.DeleteUserResponse) error {
	err := userDataService.userDataClient.DeleteUser(request.UserID, request.Username, request.Password)
	if err != nil {
		if err != config.InvalidAuthDataError {
			serviceLog.Criticalf("DeleteUser With UserID %s Returned Error %s", request.UserID, err.Error())
		}
		response.Error = err.Error()
	}
	return err
}
