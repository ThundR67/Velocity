package clients

import (
	"context"

	"github.com/SonicRoshan/Velocity/global/config"
	proto "github.com/SonicRoshan/Velocity/services/users-srv/proto"
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
	mainData config.UserMain, extraData config.UserExtra) (string, string, error) {

	mainDataProto := proto.UserMain{}
	extraDataProto := proto.UserExtra{}

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

//Get is used to get user  main data
func (usersClient UsersClient) Get(userID string) (config.UserMain, string, error) {
	request := proto.GetRequest{
		UserID: userID,
	}

	response, err := usersClient.client.Get(context.TODO(), &request)
	if err != nil {
		err = errors.Wrap(err, "Error while getting user main data through client")
		return config.UserMain{}, "", err
	} else if response.Message != "" {
		return config.UserMain{}, response.Message, nil
	}

	var data config.UserMain
	err = copy(&data, &response)
	return data, "", err
}

//Update is used to update user data
func (usersClient UsersClient) Update(userID string, update config.UserMain) error {
	var updateRequest proto.UserMain
	err := copy(&updateRequest, &update)

	if err != nil {
		return err
	}

	request := proto.UpdateRequest{
		UserID: userID,
		Update: &updateRequest,
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

//Activate is used to mark an account active
func (usersClient UsersClient) Activate(email string) (string, error) {
	request := proto.ActivateRequest{
		Email: email,
	}
	response, err := usersClient.client.Activate(context.TODO(), &request)
	if err != nil {
		return "", errors.Wrap(err, "Error While Activating User Through Client")
	}
	return response.Message, nil
}
