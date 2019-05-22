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

func (userDataService *UserDataService) getData(ctx context.Context, request *proto.GetUserRequest, collection string) (map[string]interface{}, error) {
	userData, err := userDataService.userDataClient.GetUserData(request.UserID, collection)
	if err != nil {
		serviceLog.Warnf("getData With Collection %s With UserID %s Returned Error %+v", collection, request.UserID, err)
		return nil, err
	}
	return userData, err
}

func (userDataService *UserDataService) updateData(ctx context.Context, request *proto.UpdateUserRequest, collection string) error {
	err := userDataService.userDataClient.UpdateUserData(request.UserID, request.Field, request.UpdatedValue, collection)
	if err != nil {
		serviceLog.Warnf("UpdateUser With Id %s Returned Error %+v", request.UserID, err)
		return err
	}
	return err
}

//AddUser Handles AddUser Function
func (userDataService UserDataService) AddUser(ctx context.Context, request *proto.UserData, response *proto.AddUserResponse) error {
	userMainData := castFromStringToInterface(request.UserMainData)
	userExtraData := castFromStringToInterface(request.UserExtraData)
	userID, msg, err := userDataService.userDataClient.AddUser(userMainData, userExtraData)
	if err != nil {
		serviceLog.Warnf("AddUser With Data %v Returned Error %+v", userMainData, err)
		return err
	}
	response.Message = msg
	response.UserID = userID
	return err
}

//GetUserData Handles GetUserData Function
func (userDataService UserDataService) GetUserData(ctx context.Context, request *proto.GetUserRequest, response *proto.UserData) error {
	data, err := userDataService.getData(ctx, request, config.DBConfigUserDataCollection)
	response.UserMainData = castFromInterfaceToString(data)
	return err
}

//GetUserExtraData Handles GetUserExtraData Function
func (userDataService UserDataService) GetUserExtraData(ctx context.Context, request *proto.GetUserRequest, response *proto.UserData) error {
	data, err := userDataService.getData(ctx, request, config.DBConfigUserExtraDataCollection)
	response.UserExtraData = castFromInterfaceToString(data)
	return err
}

//GetUserByUsernameOrEmail Handles GetUserByUsernameOrEmail Function
func (userDataService UserDataService) GetUserByUsernameOrEmail(ctx context.Context, request *proto.GetUserByUsernameOrEmailRequest, response *proto.UserData) error {
	userData, msg, err := userDataService.userDataClient.GetUserByUsernameOrEmail(request.Username, request.Email)
	if err != nil {
		serviceLog.Warnf("GetUserByUsernameOrEmail With Username %s Email %s Returned Error %+v", request.Username, request.Email, err)
		return err
	}
	response.Message = msg
	response.UserMainData = castFromInterfaceToString(userData)
	return err
}

//AuthUser Handles AuthUser Function
func (userDataService UserDataService) AuthUser(ctx context.Context, request *proto.AuthUserRequest, response *proto.AuthUserResponse) error {
	valid, userID, msg, err := userDataService.userDataClient.AuthUser(request.Username, request.Email, request.Password)
	if err != nil {
		serviceLog.Warnf("AuthUser With Username %s Email %s Returned Error %+v", request.Username, request.Email, err)
		return err
	}
	response.Message = msg
	response.Valid = valid
	response.UserID = userID
	return err
}

//UpdateUserData Handles UpdateUserData Function
func (userDataService UserDataService) UpdateUserData(ctx context.Context, request *proto.UpdateUserRequest, response *proto.UpdateUserResponse) error {
	return userDataService.updateData(ctx, request, config.DBConfigUserDataCollection)
}

//UpdateUserExtraData Handles UpdateUserData Function
func (userDataService UserDataService) UpdateUserExtraData(ctx context.Context, request *proto.UpdateUserRequest, response *proto.UpdateUserResponse) error {
	return userDataService.updateData(ctx, request, config.DBConfigUserExtraDataCollection)
}

//DeleteUser Handles DeleteUser Function
func (userDataService UserDataService) DeleteUser(ctx context.Context, request *proto.DeleteUserRequest, response *proto.DeleteUserResponse) error {
	msg, err := userDataService.userDataClient.DeleteUser(request.UserID, request.Username, request.Password)
	if err != nil {
		serviceLog.Warnf("DeleteUser With UserID %s Returned Error %+v", request.UserID, err)
		return err
	}
	response.Message = msg
	return err
}
