package clients

import (
	"context"

	"github.com/SonicRoshan/Velocity/global/config"
	proto "github.com/SonicRoshan/Velocity/services/users-srv/proto"
	"github.com/gorilla/schema"
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
	mainData config.UserMain, extraData config.UserExtra) (string, string) {

	mainDataProto := proto.UserMain{}
	extraDataProto := proto.UserExtra{}

	copy(&mainDataProto, &mainData)
	copy(&extraDataProto, &extraData)

	request := proto.AddRequest{
		MainData:  &mainDataProto,
		ExtraData: &extraDataProto,
	}

	response, _ := usersClient.client.Add(context.TODO(), &request)

	return response.UserID, response.Message
}

//Get is used to get user  main data
func (usersClient UsersClient) Get(userID string) (config.UserMain, error) {
	request := proto.GetRequest{
		UserID: userID,
	}

	response, err := usersClient.client.Get(context.TODO(), &request)
	if err != nil {
		err = errors.Wrap(err, "Error while getting user main data through client")
		return config.UserMain{}, err
	}

	var data config.UserMain
	err = copy(&data, &response)
	return data, err
}

//GetExtra is used to get user extra data
func (usersClient UsersClient) GetExtra(userID string) (config.UserExtra, error) {
	request := proto.GetRequest{
		UserID: userID,
	}

	response, err := usersClient.client.GetExtra(context.TODO(), &request)
	if err != nil {
		err = errors.Wrap(err, "Error while getting user extra data through client")
		return config.UserExtra{}, err
	}
	var data config.UserExtra
	err = copy(&data, &response)
	return data, err
}

//GetByUsernameOrEmail is used to get a user based on username or email
func (usersClient UsersClient) GetByUsernameOrEmail(
	username, email string) (config.UserMain, string) {

	request := proto.GetByUsernameOrEmailRequest{
		Username: username,
		Email:    email,
	}

	response, _ := usersClient.client.GetByUsernameOrEmail(context.TODO(), &request)

	var data config.UserMain
	copy(&data, &response)
	return data, response.Message
}

//Update is used to update user data
func (usersClient UsersClient) Update(userID string, update config.UserMain) {
	var updateRequest proto.UserMain
	copy(&updateRequest, &update)

	request := proto.UpdateRequest{
		UserID: userID,
		Update: &updateRequest,
	}

	usersClient.client.Update(context.TODO(), &request)
}

//UpdateExtra is used to update user's extra data
func (usersClient UsersClient) UpdateExtra(userID string, update config.UserExtra) {
	var updateRequest proto.UserExtra
	copy(&updateRequest, &update)

	request := proto.UpdateExtraRequest{
		UserID: userID,
		Update: &updateRequest,
	}

	usersClient.client.UpdateExtra(context.TODO(), &request)
}

//Auth is used to authenticate a user
func (usersClient UsersClient) Auth(username, password string) (string, string, error) {
	request := proto.AuthRequest{
		Username: username,
		Password: password,
	}
	response, _ := usersClient.client.Auth(context.TODO(), &request)
	return response.UserID, response.Message, nil
}

//Activate is used to mark an account active
func (usersClient UsersClient) Activate(email string) string {
	request := proto.ActivateRequest{
		Email: email,
	}
	response, _ := usersClient.client.Activate(context.TODO(), &request)
	return response.Message
}
