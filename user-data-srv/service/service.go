package userdataservice

import (
	"context"

	"github.com/SonicRoshan/Velocity/user-data-srv/config"
	proto "github.com/SonicRoshan/Velocity/user-data-srv/proto"
	userDataManager "github.com/SonicRoshan/Velocity/user-data-srv/user-data-manager"
	logger "github.com/jex-lin/golang-logger"
)

//Log Loding Logger
var Log = logger.NewLogFile("logs/service_logs.log")

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
		Log.Criticalf("AddUser With Data %v Returned Error %s", request.Data, err)
	} else if userID == config.UsernameExistsMsg || userID == config.EmailExistsMsg {
		response.Error = userID
	} else {
		response.UserID = userID
	}
	return err
}

//GetUser Handles GetUser Function
func (userDataService UserDataService) GetUser(ctx context.Context, request *proto.GetUserRequest, response *proto.UserData) error {
	userData, err := userDataService.userDataClient.GetUser(request.UserID)
	if err != nil {
		Log.Criticalf("GetUser With UserID %s Returned Error %s", request.UserID, err)
		response.Error = err.Error()
	}
	response.Data = castFromInterfaceToString(userData)
	return err
}

//GetUserByUsernameOrEmail Handles GetUserByUsernameOrEmail Function
func (userDataService UserDataService) GetUserByUsernameOrEmail(ctx context.Context, request *proto.GetUserByUsernameOrEmailRequest, response *proto.UserData) error {
	userData, err := userDataService.userDataClient.GetUserByUsernameOrEmail(request.Username, request.Email)
	if err != nil {
		Log.Criticalf("GetUserByUsernameOrEmail With Username %s Email %s Returned Error %s", request.Username, request.Email, err)
		response.Error = err.Error()
	}
	response.Data = castFromInterfaceToString(userData)
	return err
}

//AuthUser Handles AuthUser Function
func (userDataService UserDataService) AuthUser(ctx context.Context, request *proto.AuthUserRequest, response *proto.AuthUserResponse) error {
	valid, userID, err := userDataService.userDataClient.AuthUser(request.Username, request.Email, request.Password)
	if err != nil {
		Log.Criticalf("AuthUser With Username %s Email %s Returned Error %s", request.Username, request.Email, err)
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
		Log.Criticalf("UpdateUser With Id %s Returned Error %s", request.UserID, err)
		response.Error = err.Error()
	}
	return err

}

//DeleteUser Handles DeleteUser Function
func (userDataService UserDataService) DeleteUser(ctx context.Context, request *proto.DeleteUserRequest, response *proto.DeleteUserResponse) error {
	err := userDataService.userDataClient.DeleteUser(request.UserID, request.Username, request.Password)
	if err != nil {
		Log.Criticalf("DeleteUser With UserID %s Returned Error %s", request.UserID, err.Error())
		response.Error = err.Error()
	}
	return err
}
