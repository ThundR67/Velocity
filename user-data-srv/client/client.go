package userdatasrvclient

import (
	userDataSrvProto "github.com/SonicRoshan/Velocity/user-data-srv/proto"
	micro "github.com/micro/go-micro"
)

//GetUserDataSrvClient returns a client for user-data-srv
func GetUserDataSrvClient(service micro.Service) userDataSrvProto.UserDataManagerService {
	return userDataSrvProto.NewUserDataManagerService("user-data-srv", service.Client())
}
