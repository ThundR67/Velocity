package jwtsrvclient

import (
	proto "github.com/SonicRoshan/Velocity/jwt-srv/proto"
	micro "github.com/micro/go-micro"
)

//GetJWTClient returns a client for JWT Service
func GetJWTClient(service micro.Service) proto.JWTManagerService {
	return proto.NewJWTManagerService("jwt-srv", service.Client())
}
