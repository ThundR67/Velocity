package users

import (
	"github.com/SonicRoshan/global/clients"
	micro "github.com/micro/go-micro"
)

//Resource is the API For User Data
type Resource struct {
    jwt clients.JWTClient
}

//Init initializes
func (usersResource *UsersResource) Init(service micro.Service) {
    usersResource.jwt = clients.NewJWTClient(authSrv)
}

//FindUser is used to find a user by username or email
func (usersResource UsersResource) FindUser(params graphql.params) (interface{}, error) {
	return 32, nil
}