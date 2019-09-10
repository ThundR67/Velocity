package users

import (
	"github.com/SonicRoshan/straf"
	"github.com/SonicRoshan/global/config"
	micro "github.com/micro/go-micro"
)


func GetResource(service micro.Service) {
	userType, err := straf.GetGraphQLObject(config.User{})
}