package users

import (
	"strings"

	"github.com/SonicRoshan/scope"
	"github.com/graphql-go/graphql"
	micro "github.com/micro/go-micro"

	"github.com/SonicRoshan/Velocity/global/clients"
	"github.com/SonicRoshan/Velocity/global/config"
	"github.com/SonicRoshan/Velocity/global/utils"
)

//Resource is the API For User Data
type resource struct {
	users clients.UsersClient
}

//Init initializes
func (resource *resource) Init(service micro.Service) {
	resource.users = clients.NewUsersClient(service)
}

func (resource resource) middleware(
	function func(graphql.ResolveParams) (interface{}, error),
	params graphql.ResolveParams) (interface{}, error) {

	output, err := function(params)
	userOutput := output.(UserResponse)
	if err != nil {
		return nil, err
	}

	scopesHeld := params.Context.Value(config.ResourceSrvConfigScopesKey).([]string)
	userID := params.Context.Value(config.ResourceSrvConfigUserIDKey)
	userIDOnFunction := userOutput.UserMain.UserID
	if userIDOnFunction == "" {
		userIDOnFunction = userOutput.UserExtra.UserID
	}

	/*
		Here,
		    if userID is equal to userID on which function will run,
		Then,
		    from scopes all scopes starting with @me:user will be added to filteredScopes
		Else,
		    all scopes from scopes starting with @all:user will be added to filteredScopes
	*/

	toCheck := config.ScopesConfigAllTag
	if userID == userIDOnFunction {
		toCheck = config.ScopesConfigMeTag
	}

	filteredScopes := []string{}

	for _, currentScope := range scopesHeld {
		scopeSplit := strings.Split(currentScope, config.ScopesConfigDivider)
		if scopeSplit[0] == toCheck && scopeSplit[1] == config.ScopeConfigUserScope {
			//Removing entity tag from scope
			currentScope = strings.Join(scopeSplit[1:], config.ScopesConfigDivider)
			filteredScopes = append(filteredScopes, currentScope)
		}
	}

	scope.FilterRead(&userOutput.UserMain, filteredScopes)
	scope.FilterRead(&userOutput.UserExtra, filteredScopes)

	return userOutput, nil
}

//getUser is used to find a user by username or email
func (resource resource) getUser(params graphql.ResolveParams) (interface{}, error) {
	var userParams config.UserMain
	var response UserResponse
	utils.UnmarshalMap(params.Args, &userParams)

	if userParams.UserID != "" {
		user, err := resource.users.Get(userParams.UserID)
		response = UserResponse{UserMain: user, Message: ""}
		return response, err

	} else if userParams.Username != "" || userParams.Email != "" {
		user, msg := resource.users.GetByUsernameOrEmail(
			userParams.Username,
			userParams.Email,
		)
		if msg != "" {
			response = UserResponse{UserMain: user, Message: msg}
			return response, nil
		}

		response = UserResponse{UserMain: user, Message: msg}
		return response, nil
	}

	return UserResponse{Message: config.InvalidInputMsg}, nil
}

//getUserExtra is used to get a user extra data
func (resource resource) getUserExtra(params graphql.ResolveParams) (interface{}, error) {
	var userParams config.UserMain
	utils.UnmarshalMap(params.Args, &userParams)
	data, err := resource.users.GetExtra(userParams.UserID)
	response := UserResponse{UserExtra: data}
	return response, err
}

//updateUser is used to update main data of the user
func (resource resource) updateUser(params graphql.ResolveParams) (interface{}, error) {
	var update config.UserMain
	utils.UnmarshalMap(params.Args, &update)
	userID := update.UserID
	update.UserID = ""
	resource.users.Update(userID, update)
	return nil, nil
}

//updateUserExtra is used to update extra data of the user
func (resource resource) updateUserExtra(params graphql.ResolveParams) (interface{}, error) {
	var update config.UserExtra
	utils.UnmarshalMap(params.Args, &update)
	userID := update.UserID
	update.UserID = ""
	resource.users.UpdateExtra(userID, update)
	return nil, nil
}
