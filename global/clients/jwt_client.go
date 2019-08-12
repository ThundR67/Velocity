package clients

import (
	"context"

	"github.com/SonicRoshan/Velocity/global/config"
	proto "github.com/SonicRoshan/Velocity/services/jwt-srv/proto"
	micro "github.com/micro/go-micro"
	"github.com/pkg/errors"
)

//NewJWTClient is used to make a JWT Client
func NewJWTClient(service micro.Service) JWTClient {
	client := JWTClient{}
	client.Init(service)
	return client
}

//JWTClient is jwt service client
type JWTClient struct {
	client proto.JWTService
}

//Init initializes cliet
func (jwtClient *JWTClient) Init(service micro.Service) {
	jwtClient.client = proto.NewJWTService(config.JWTService, service.Client())
}

//FreshToken is used to create a fresh access token
func (jwtClient JWTClient) FreshToken(userIdentity string) (string, error) {
	request := proto.JWTData{
		UserIdentity: userIdentity,
	}
	response, err := jwtClient.client.FreshToken(context.TODO(), &request)
	if err != nil {
		return "", errors.Wrap(
			err, "Error While Generating Fresh Access Token Through Client Through Client")
	}
	return response.Token, err
}

//AccessAndRefreshTokens is used create access and refresh tokens
func (jwtClient JWTClient) AccessAndRefreshTokens(
	userIdentity string, scopes []string) (string, string, string, error) {

	request := proto.JWTData{
		UserIdentity: userIdentity,
		Scopes:       scopes,
	}
	response, err := jwtClient.client.AccessAndRefreshTokens(context.TODO(), &request)
	if err != nil {
		err = errors.Wrap(
			err, "Error While Generating Access And Refresh Token Through Client")
		return "", "", "", err
	}
	return response.AcccessToken, response.RefreshToken, response.Message, nil
}

//RefreshTokens is used to create access and refresh token based on previous refresh token
func (jwtClient JWTClient) RefreshTokens(
	refreshToken string) (string, string, string, error) {

	request := proto.Token{
		Token: refreshToken,
	}
	response, err := jwtClient.client.RefreshTokens(context.TODO(), &request)
	if err != nil {
		err = errors.Wrap(
			err, `Error While Generating Access 
				  And Refresh Token Bases On Refresh Token Through Client`)

		return "", "", "", err
	}
	return response.AcccessToken, response.RefreshToken, response.Message, nil
}

//ValidateToken is used to validate a token
func (jwtClient JWTClient) ValidateToken(
	token, tokenType string) (bool, string, []string, error) {

	request := proto.Token{
		Token:     token,
		TokenType: tokenType,
	}

	response, err := jwtClient.client.ValidateToken(context.TODO(), &request)
	if err != nil {
		err = errors.Wrap(err, "Error While Validating Token Through Client")
		return false, "", nil, err
	}
	userIdentity := response.UserIdentity
	scopes := response.Scopes
	return response.Valid, userIdentity, scopes, nil
}
