package users

import (
	"github.com/SonicRoshan/Velocity/global/config"
)

//UserResponse is used to send response from graphQL server
type UserResponse struct {
	UserMain  config.UserMain
	UserExtra config.UserExtra
	Message   string
}
