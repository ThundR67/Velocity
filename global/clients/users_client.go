package clients

import (
	"context"

	"github.com/SonicRoshan/Velocity/global/config"
	proto "github.com/SonicRoshan/Velocity/users-srv/proto"
	"github.com/jinzhu/copier"
	micro "github.com/micro/go-micro"
	"github.com/pkg/errors"
)

//NewUsersClient is used to get users client struct
func NewUsersClient(service micro.Service) UsersClient {
	usersClient := UsersClient{}
	usersClient.Init(service)
	return usersClient
}

func copy(to interface{}, from interface{}) error {
	err := copier.Copy(to, from)
	if err != nil {
		return errors.Wrap(err, "Error While Copying")
	}
	return nil
}

//UsersClient is used to communicate with users service
type UsersClient struct {
	client proto.UsersService
}

//Init is used to initialise client
func (usersClient *UsersClient) Init(service micro.Service) {
	usersClient.client = proto.NewUsersService(config.UsersService, service.Client())
}

//Add is used to add a user
func (usersClient UsersClient) Add(
	mainData, extraData config.UserType) (string, string, error) {

	mainDataProto := proto.User{}
	extraDataProto := proto.User{}

	err := copy(&mainDataProto, &mainData)
	if err != nil {
		return "", "", err
	}
	err = copy(&extraDataProto, &extraData)
	if err != nil {
		return "", "", err
	}

	request := proto.AddRequest{
		MainData:  &mainDataProto,
		ExtraData: &extraDataProto,
	}

	response, err := usersClient.client.Add(context.TODO(), &request)
	if err != nil {
		err = errors.Wrap(err, "Error while adding user through client")
		return "", "", err
	}

	return response.UserID, response.Message, nil
}

//Get is used to get user data inn a collection
func (usersClient UsersClient) Get(userID, collection string) (config.UserType, string, error) {
	request := proto.GetRequest{
		UserID:     userID,
		Collection: collection,
	}

	response, err := usersClient.client.Get(context.TODO(), &request)
	if err != nil {
		err = errors.Wrap(err, "Error while getting user data through client")
		return config.UserType{}, "", err
	} else if response.Message != "" {
		return config.UserType{}, response.Message, nil
	}

	var data config.UserType
	err = copy(&data, &response)
	return data, "", err
}

//Update is used to update user data
func (usersClient UsersClient) Update(userID, collection string, update config.UserType) error {
	var updateRequest proto.User
	err := copy(&updateRequest, &update)

	if err != nil {
		return nil
	}

	request := proto.UpdateRequest{
		UserID:     userID,
		Collection: collection,
		Update:     &updateRequest,
	}

	_, err = usersClient.client.Update(context.TODO(), &request)
	if err != nil {
		return errors.Wrap(err, "Error while updating through client")
	}
	return nil
}

//Auth is used to authenticate a user
func (usersClient UsersClient) Auth(username, password string) (string, string, error) {
	request := proto.AuthRequest{
		Username: username,
		Password: password,
	}
	response, err := usersClient.client.Auth(context.TODO(), &request)
	if err != nil {
		errors.Wrap(err, "Error while authenticating user through client")
		return "", "", err
	}
	return response.UserID, response.Message, nil
}
