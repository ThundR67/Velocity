package users

import (
	"github.com/SonicRoshan/straf"
	"github.com/graphql-go/graphql"
	micro "github.com/micro/go-micro"

	"github.com/SonicRoshan/Velocity/global/config"
)

//GetResource returns graphQL Fields
func GetResource(service micro.Service) (graphql.Fields, error) {
	responseType, err := straf.GetGraphQLObject(UserResponse{})
	if err != nil {
		return graphql.Fields{}, err
	}

	resource := resource{}
	resource.Init(service)

	builder := straf.NewSchemaBuilder(responseType, config.UserMain{}, resource.middleware)

	builder.AddFunction(
		"GetUser",
		"This function is used to get users main data",
		resource.getUser,
	)

	builder.AddFunction(
		"GetUserExtra",
		"This function is used to get users extra data",
		resource.getUserExtra,
	)

	builder.AddFunction(
		"UpdateUser",
		"This function is used to update users main data",
		resource.updateUser,
	)

	builder.AddFunction(
		"UpdateUserExtra",
		"This function is used to update users extra data",
		resource.updateUserExtra,
	)

	return builder.Schema, nil
}
