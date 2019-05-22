package clients

import (
	"context"

	userDataSrvProto "github.com/SonicRoshan/Velocity/user-data-srv/proto"
	micro "github.com/micro/go-micro"
	"github.com/pkg/errors"
)

//GetUserDataSrvClient returns a client for user-data-srv
func GetUserDataSrvClient(service micro.Service) userDataSrvProto.UserDataManagerService {
	return userDataSrvProto.NewUserDataManagerService("user-data-srv", service.Client())
}

//UserDataClient is a client to user data service
type UserDataClient struct {
	client userDataSrvProto.UserDataManagerService
}

//Init initialises client
func (userDataClient *UserDataClient) Init(service micro.Service) {
	userDataClient.client = userDataSrvProto.NewUserDataManagerService("user-data-srv", service.Client())
}

//AddUser adds an user
func (userDataClient UserDataClient) AddUser(userMainData, userExtraData map[string]string) (string, string, error) {
	request := userDataSrvProto.UserData{
		UserMainData:  userMainData,
		UserExtraData: userExtraData,
	}
	response, err := userDataClient.client.AddUser(context.TODO(), &request)
	if err != nil {
		errors.Wrap(err, "Error while adding user through client")
		return "", "", err
	}
	return response.UserID, response.Message, nil
}

//AuthUser auths user
func (userDataClient UserDataClient) AuthUser(username, password string) (string, string, error) {
	request := userDataSrvProto.AuthUserRequest{
		Username: username,
		Password: password,
	}
	response, err := userDataClient.client.AuthUser(context.TODO(), &request)
	if err != nil {
		errors.Wrap(err, "Error while auth user through client")
		return "", "", err
	}
	return response.UserID, response.Message, nil
}
